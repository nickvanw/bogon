package state

import (
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type InviteHandler struct {
	Bot   *ircx.Bot
	State *State
}

func (h *InviteHandler) Handle(s irc.Sender, m *irc.Message) {
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
