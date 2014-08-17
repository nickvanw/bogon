package state

import (
	"fmt"
	"log"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type ModeHandler struct {
	Bot   *ircx.Bot
	State *State
}

func (h *ModeHandler) Handle(s irc.Sender, m *irc.Message) {
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
	for _, v := range modeString {
		switch v {
		case '+':
			plus = true
		case '-':
			plus = false
		case 'q', 'a', 'o', 'h', 'v':
			nick := modeArgs[0]
			user := channel.GetUser(nick)
			if plus {
				fmt.Println("Adding mode", v, "to user", nick)
				user.Modes[v] = struct{}{}
			} else {
				fmt.Println("removing mode", v, "to user", nick)
				delete(user.Modes, v)
			}
			modeArgs = modeArgs[1:]
		default:
			if plus {
				fmt.Println("Adding Channel Mode", v)
				channel.Modes[v] = struct{}{}
			} else {
				fmt.Println("Removing Mode", v)
				delete(channel.Modes, v)
			}
		}
	}
}
