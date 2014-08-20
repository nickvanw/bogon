package cmd

import (
	"fmt"

	"github.com/nickvanw/bogon/state"
)

func init() {
	AddPlugin("ChanStat", "(?i)^\\.chanstat(s)?$", MessageHandler(ChanStat), false, true)
}

func ChanStat(msg *Message) {
	var ch *state.Channel
	var err error
	if len(msg.Params) > 1 {
		ch, err = msg.State.GetChan(msg.Params[1])
	} else {
		ch, err = msg.State.GetChan(msg.To)
	}
	if err != nil {
		msg.Return("I couldn't even find this channel, I'm not working well")
		return
	}
	numPeople := len(ch.Users)
	ops := ch.Ops()
	voice := ch.Voice()
	msg.Return(fmt.Sprintf("%s: %d people, %d ops, %d voice", ch.Name, numPeople, ops, voice))
}
