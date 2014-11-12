package state

import (
	"log"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type ModeHandler struct {
	Bot   *ircx.Bot
	State *State
}

func (h *ModeHandler) Handle(s ircx.Sender, m *irc.Message) {
	h.State.ParseModes(m.Params)
}

func (s *State) ParseModes(modes []string) {
	channel, err := s.GetChan(modes[0])
	if err != nil {
		log.Println("I got modes for a channel I am not in")
		return
	}
	modeString := modes[1]
	modeArgs := modes[2:]
	var plus bool
	s.Lock()
	defer s.Unlock()
	for _, v := range modeString {
		switch v {
		case '+':
			plus = true
		case '-':
			plus = false
		case 'q', 'a', 'o', 'h', 'v':
			nick := modeArgs[0]
			user, err := channel.GetUser(nick)
			if err != nil {
				continue
			}
			if plus {
				user.Modes[v] = struct{}{}
			} else {
				delete(user.Modes, v)
			}
			modeArgs = modeArgs[1:]
		default:
			if plus {
				channel.Modes[v] = struct{}{}
			} else {
				delete(channel.Modes, v)
			}
		}
	}
}
