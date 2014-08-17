package cmd

import "fmt"

func init() {
	AddPlugin("ChanStat", "(?i)^\\.chanstat(s)?$", MessageHandler(ChanStat), false, false)
}
func ChanStat(msg *Message) {
	ch, err := msg.State.GetChan(msg.To)
	msg.Return(fmt.Sprintf("%s, %s", ch, err))
}
