package state

import (
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type JoinHandler struct {
	Bot   *ircx.Bot
	State *State
}

func (h *JoinHandler) Handle(s irc.Sender, m *irc.Message) {
	var channel string
	if len(m.Params) > 0 {
		channel = m.Params[0]
	} else {
		channel = m.Trailing
	}
	if m.Prefix.Name == h.State.Name {
		h.State.NewChannel(channel)
	} else {
		h.State.NewUser(channel, m.Prefix.Name)
	}
}
