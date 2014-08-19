package cmd

func init() {
	AddPlugin("EndEncryption", "(?i)^\\.encoff?$", MessageHandler(EndEncryption), false, false)
}

func EndEncryption(msg *Message) {
	if _, ok := msg.State.Encryption[msg.To]; ok {
		delete(msg.State.Encryption, msg.To)
		msg.Return("I've removed your encryption key")
	}
}
