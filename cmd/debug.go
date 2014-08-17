package cmd

import "fmt"

func init() {
	AddPlugin("ChanStat", "(?i)^\\.chanstat(s)?$", MessageHandler(ChanStat), false, false)
}

func ChanStat(msg *Message) {
	ch, err := msg.State.GetChan(msg.To)
	if err != nil {
		msg.Return("I couldn't even find this channel, I'm not working well")
		return
	}
	numPeople := len(ch.Users)
	ops := ch.Ops()
	voice := ch.Voice()
	msg.Return(fmt.Sprintf("%s: %d people, %d ops, %d voice", ch.Name, numPeople, ops, voice))
}
