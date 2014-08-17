package cmd

import (
	"strings"

	"github.com/nickvanw/bogon/state"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type CommandHandler struct {
	Bot   *ircx.Bot
	State *state.State
}

func (cmd *CommandHandler) Handle(s irc.Sender, m *irc.Message) {
	for _, v := range Commands {
		data := strings.Split(m.Trailing, " ")
		if v.Command.MatchString(data[0]) || v.Raw {
			sendMessage := &Message{
				Params: data,
				Sender: s,
				State:  cmd.State,
			}
			// Transparently make Return send message to the channel the
			// message was sent in, or to the user in PM
			if strings.ContainsAny(m.Params[0], strings.Join([]string{string(irc.Channel), string(irc.Distributed)}, "")) {
				sendMessage.To = m.Params[0]
			} else {
				sendMessage.To = m.Name
			}
			sendMessage.Name = v.Name
			go v.Function(sendMessage)
		}
	}
}

type Message struct {
	Params []string
	Sender irc.Sender
	To     string
	State  *state.State
	Name   string
}

func (m *Message) Return(out string) {
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
