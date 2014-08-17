package state

import (
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type QuitHandler struct {
	Bot   *ircx.Bot
	State *State
}

func (h *QuitHandler) Handle(s irc.Sender, m *irc.Message) {
	h.State.QuitUser(m.Prefix.Name)
}

func (s *State) QuitUser(name string) {
	s.Lock()
	defer s.Unlock()
	for _, channel := range s.Channels {
		channel.RemoveUser(name)
	}
}
