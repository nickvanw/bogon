package bogon

import (
	"fmt"
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func (c *Client) registerStateHandlers() {
	c.bot.HandleFunc(irc.INVITE, c.inviteHandler)
	c.bot.HandleFunc(irc.JOIN, c.joinHandler)
	c.bot.HandleFunc(irc.KICK, c.kickHandler)
	c.bot.HandleFunc(irc.MODE, c.modeHandler)
	c.bot.HandleFunc(irc.RPL_NAMREPLY, c.namesHandler)
	c.bot.HandleFunc(irc.NICK, c.nickHandler)
	c.bot.HandleFunc(irc.ERR_NICKNAMEINUSE, c.nickTakenHandler)
	c.bot.HandleFunc(irc.PART, c.partHandler)
	c.bot.HandleFunc(irc.QUIT, c.quitHandler)
	c.bot.HandleFunc(irc.RPL_TOPIC, c.topicHandler)

	// register handlers for encryption
	c.bot.HandleFunc(irc.NOTICE, c.encryptionStartHandler)
	c.bot.HandleFunc(irc.PRIVMSG, c.encryptionMsgHandler)

	// register command handler
	c.bot.HandleFunc(irc.PRIVMSG, c.commandHandler)
}

func (c *Client) inviteHandler(s ircx.Sender, m *irc.Message) {
	var channel string
	if len(m.Params) > 1 {
		channel = m.Params[1]
	} else {
		channel = m.Trailing
	}
	msg := &irc.Message{
		Command: irc.JOIN,
		Params:  []string{channel},
	}
	s.Send(msg)
}

func (c *Client) joinHandler(s ircx.Sender, m *irc.Message) {
	var channel string
	if len(m.Params) > 0 {
		channel = m.Params[0]
	} else {
		channel = m.Trailing
	}
	if m.Prefix.Name == c.state.Name() {
		c.state.NewChannel(channel)
	} else {
		c.state.NewUser(channel, m.Prefix.Name)
	}
}

func (c *Client) kickHandler(s ircx.Sender, m *irc.Message) {
	if m.Params[1] == c.state.Name() {
		c.state.RemoveChannel(m.Params[0])
		if c.state.Rejoin() {
			msg := &irc.Message{
				Command: irc.JOIN,
				Params:  []string{m.Params[0]},
			}
			s.Send(msg)
		}
	} else {
		c.state.RemoveUser(m.Params[0], m.Params[1])
	}
}

func (c *Client) modeHandler(s ircx.Sender, m *irc.Message) {
	c.state.ParseModes(m.Params)
}

func (c *Client) namesHandler(s ircx.Sender, m *irc.Message) {
	channel, err := c.state.GetChan(m.Params[2])
	if err != nil {
		return
	}
	newNames := strings.Split(m.Trailing, " ")
	channel.NewUser(newNames...)
}

func (c *Client) nickHandler(s ircx.Sender, m *irc.Message) {
	if m.Name == c.state.Name() {
		c.state.SetName(m.Trailing)
	}
	c.state.RenameUser(m.Name, m.Trailing)
}

func (c *Client) nickTakenHandler(s ircx.Sender, m *irc.Message) {
	newName := fmt.Sprintf("%s|", c.state.Name())
	msg := &irc.Message{
		Command: irc.NICK,
		Params:  []string{newName},
	}
	c.state.SetName(newName)
	s.Send(msg)
}

func (c *Client) partHandler(s ircx.Sender, m *irc.Message) {
	if m.Prefix.Name == c.state.Name() {
		c.state.RemoveChannel(m.Params[0])
	} else {
		c.state.RemoveUser(m.Params[0], m.Prefix.Name)
	}
}

func (c *Client) quitHandler(s ircx.Sender, m *irc.Message) {
	c.state.QuitUser(m.Prefix.Name)
}

func (c *Client) topicHandler(s ircx.Sender, m *irc.Message) {
	channel, err := c.state.GetChan(m.Params[1])
	if err != nil {
		return
	}
	channel.SetTopic(m.Trailing)
}
