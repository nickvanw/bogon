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
	if m.Params[1] == h.State.Name {
		h.State.RemoveChannel(m.Params[0])
		if on, ok := h.Bot.Options["rejoin"]; ok && on {
			msg := &irc.Message{
				Command: irc.JOIN,
				Params:  []string{m.Params[0]},
			}
			s.Send(msg)
		}
	} else {
		h.State.RemoveUser(m.Params[0], m.Params[1])
	}
}
