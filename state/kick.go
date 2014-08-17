package state

import (
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type KickHandler struct {
	Bot   *ircx.Bot
	State *State
}

func (h *KickHandler) Handle(s irc.Sender, m *irc.Message) {
	if m.Prefix.Name == h.State.Name {
		h.State.RemoveChannel(m.Params[0])
	} else {
		h.State.RemoveUser(m.Params[0], m.Prefix.Name)
	}
}
