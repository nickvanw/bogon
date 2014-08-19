package state

import (
	"fmt"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type NickHandler struct {
	Bot   *ircx.Bot
	State *State
}

func (h *NickHandler) Handle(s irc.Sender, m *irc.Message) {
	if m.Name == h.State.Name {
		h.State.Name = m.Trailing
	}
	h.State.RenameUser(m.Name, m.Trailing)
}

func (s *State) RenameUser(oldname string, newname string) {
	s.Lock()
	for _, channel := range s.Channels {
		olduser := channel.GetUser(oldname)
		olduser.Name = newname
	}
	if data, ok := s.Encryption[oldname]; ok {
		delete(s.Encryption, oldname)
		s.Encryption[newname] = data
	}
	s.Unlock()
}

type NickTakenHandler struct {
	Bot   *ircx.Bot
	State *State
}

func (h *NickTakenHandler) Handle(s irc.Sender, m *irc.Message) {
	msg := &irc.Message{
		Command: irc.NICK,
		Params:  []string{fmt.Sprintf("%s|", h.State.Name)},
	}
	h.State.Name = fmt.Sprintf("%s|", h.State.Name)
	s.Send(msg)
}
