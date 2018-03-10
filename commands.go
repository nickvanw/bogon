package bogon

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/dh1080"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type command struct {
	name  string
	match *regexp.Regexp
	admin bool
	raw   bool
	f     commands.CommandFunc
}

// AddCommands registers a new callback command to check
func (c *Client) AddCommands(cmds ...commands.RegisterFunc) {
	for _, v := range cmds {
		name, match, f, opts := v()
		c.commands = append(c.commands, command{name: name, match: match, f: f, admin: opts.AdminOnly, raw: opts.Raw})
	}
}

// ListCommands returns all of the registerd commands in a name, by their
// name and the regular expression they are matched to
func (c *Client) ListCommands() map[string]*regexp.Regexp {
	out := make(map[string]*regexp.Regexp)
	for _, v := range c.commands {
		out[v.name] = v.match
	}
	return out
}

func (c *Client) commandHandler(s ircx.Sender, ms *irc.Message) {
	m, enc := c.decodedMessage(ms)
	params := strings.Split(m.Trailing, " ")

	var to string
	if strings.HasPrefix(m.Params[0], string(irc.Channel)) ||
		strings.HasPrefix(m.Params[0], string(irc.Distributed)) {
		to = m.Params[0]
	} else {
		to = m.Name
	}
	msg := commands.Message{
		Params: params,
		To:     to,
		From:   m.Prefix.Name,
	}

	sender := c.commandSender(s, to, enc)
	for _, v := range c.commands {
		if v.raw || v.match.MatchString(params[0]) {
			if !v.raw {
				level.Debug(c.bot.Logger()).Log("action", "command", "command", v.name)
			}
			go noPanicCommand(msg, sender, v, c.bot.Logger())
		}
	}
}

func noPanicCommand(m commands.Message, r func(string, ...interface{}), c command, l log.Logger) {
	defer func(c command, l log.Logger) {
		if err := recover(); err != nil {
			level.Error(l).Log("action", "recover", "error", err, "command", c.name)
		}
	}(c, l)
	r(c.f(m, r))
}

func (c *Client) decodedMessage(m *irc.Message) (*irc.Message, bool) {
	if !strings.HasPrefix(m.Trailing, "+OK *") {
		return m, false
	}
	// create a copy of the message
	newMessage := *m

	if key, ok := c.state.Encryption().Check(m.Name); ok {
		decodeData := strings.TrimLeft(m.Trailing, "+OK *")
		rawData, err := base64.StdEncoding.DecodeString(decodeData)
		if err != nil {
			return m, false
		}
		data, err := dh1080.Dec(rawData, []byte(key))
		if err != nil {
			return m, false
		}
		newMessage.Trailing = string(data)

	}
	return &newMessage, true
}

func (c *Client) commandSender(sender ircx.Sender, to string, enc bool) func(s string, a ...interface{}) {
	return func(s string, a ...interface{}) {
		if s == "" {
			return
		}
		msg := s
		if len(a) > 0 {
			msg = fmt.Sprintf(msg, a...)
		}
		if enc {
			if key, ok := c.state.Encryption().Check(to); ok {
				data, err := dh1080.Enc([]byte(msg), []byte(key))
				if err == nil {
					msg = "+OK *" + base64.StdEncoding.EncodeToString(data)
				}
			}
		}
		newMsg := &irc.Message{
			Command:  irc.PRIVMSG,
			Params:   []string{to},
			Trailing: msg,
		}
		sender.Send(newMsg)
	}
}
