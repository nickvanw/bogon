package state

import (
	"log"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type TopicHandler struct {
	Bot   *ircx.Bot
	State *State
}

func (h *TopicHandler) Handle(s irc.Sender, m *irc.Message) {
	channel, err := h.State.GetChan(m.Params[1])
	if err != nil {
		log.Println("Got a topic for a channel I'm not in")
	}
	h.State.Lock()
	channel.Topic = m.Trailing
	h.State.Unlock()
}
