package state

import (
	"log"
	"strings"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type NamesHandler struct {
	Bot   *ircx.Bot
	State *State
}

func (h *NamesHandler) Handle(s irc.Sender, m *irc.Message) {
	channel, err := h.State.GetChan(m.Params[2])
	if err != nil {
		log.Println("I got names for a channel I'm not in")
	}
	newNames := strings.Split(m.Trailing, " ")
	h.State.Lock()
	for _, name := range newNames {
		channel.NewUser(name)
	}
	h.State.Unlock()
}
