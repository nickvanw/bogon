package cmd

import (
	"encoding/base64"
	"strings"

	"github.com/nickvanw/bogon/state"
	"github.com/nickvanw/bogon/state/dh1080"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type CommandHandler struct {
	Bot   *ircx.Bot
	State *state.State
}

func (cmd *CommandHandler) Handle(s irc.Sender, m *irc.Message) {
	admin := false
	if cmd.State.Admin != "" && cmd.State.Password != "" { // if we have an empty password, we NEVER want to auth anyone
		admin = cmd.State.Admin == m.Prefix.String()
	}
	enc := false
	var to, message string
	if strings.ContainsAny(m.Params[0], strings.Join([]string{string(irc.Channel), string(irc.Distributed)}, "")) {
		message = m.Trailing
		to = m.Params[0]
	} else {
		to = m.Name
		if strings.HasPrefix(m.Trailing, "+OK *") {
			if key, ok := cmd.State.Encryption[to]; ok {
				decodeData := strings.TrimLeft(m.Trailing, "+OK *")
				rawData, err := base64.StdEncoding.DecodeString(decodeData)
				if err != nil {
					return
				}
				data, err := dh1080.Dec(rawData, []byte(key))
				if err != nil {
					return
				} else {
					message = string(data)
					enc = true
				}
			}
		} else {
			message = m.Trailing
		}
	}
	for _, v := range Commands {
		data := strings.Split(message, " ")
		if v.Command.MatchString(data[0]) || v.Raw {
			if v.RequireAdmin {
				if !admin {
					continue
				}
			}
			sendMessage := &Message{
				Params:  data,
				Sender:  s,
				State:   cmd.State,
				Enc:     enc,
				To:      to,
				User:    m.Prefix,
				IsAdmin: admin,
			}
			sendMessage.Name = v.Name
			go v.Function(sendMessage)
		}
	}
}

type Message struct {
	Params  []string
	Sender  irc.Sender
	To      string
	State   *state.State
	Name    string
	Enc     bool
	User    *irc.Prefix
	IsAdmin bool
}

func (m *Message) Return(out string) {
	if m.Enc {
		if key, ok := m.State.Encryption[m.To]; ok {
			data, err := dh1080.Enc([]byte(out), []byte(key))
			if err == nil {
				out = "+OK *" + base64.StdEncoding.EncodeToString(data)
			}
		}
	}
	newMsg := &irc.Message{
		Command:  irc.PRIVMSG,
		Params:   []string{m.To},
		Trailing: out,
	}
	m.Sender.Send(newMsg)
}

type MessageHandler func(*Message)

func (f MessageHandler) Handle(m *Message) {
	f(m)
}
