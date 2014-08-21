package cmd

import "fmt"

func init() {
	AddPlugin("ChanStat", "(?i)^\\.chanstat(s)?$", MessageHandler(ChanStat), false, true)
}

func ChanStat(msg *Message) {
	var channel string
	if len(msg.Params) > 1 {
		channel = msg.Params[1]
	} else {
		channel = msg.To
	}
	ch, err := msg.State.GetChan(channel)
	if err != nil {
		msg.Return("I couldn't even find this channel, I'm not working well")
		return
	}
	numPeople := len(ch.Users)
	ops := ch.Ops()
	voice := ch.Voice()
	msg.Return(fmt.Sprintf("%s: %d people, %d ops, %d voice", channel, numPeople, ops, voice))
}
