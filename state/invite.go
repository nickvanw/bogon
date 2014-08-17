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
	msg := &irc.Message{
		Command: irc.JOIN,
		Params:  []string{m.Params[1]},
	}
	s.Send(msg)
}
