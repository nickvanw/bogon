package state

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/nickvanw/bogon/state/dh1080"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type EncryptionHandler struct {
	Bot   *ircx.Bot
	State *State
}

func (h *EncryptionHandler) Handle(s irc.Sender, m *irc.Message) {
	data := strings.Split(m.Trailing, " ")
	if len(data) < 2 {
		return
	}
	if data[0] == "DH1080_INIT" {
		enc := dh1080.DH1080_Init()
		err := enc.Unpack(m.Trailing)
		data, err := enc.Pack()
		if err != nil {
			return
		}
		msg := &irc.Message{
			Command:  irc.NOTICE,
			Params:   []string{m.Name},
			Trailing: data + " CBC",
		}
		s.Send(msg)
		secret, err := enc.GetSecret()
		if err != nil {
			msg := &irc.Message{
				Command:  irc.PRIVMSG,
				Params:   []string{m.Name},
				Trailing: "I couldn't create the key for our chat, I won't be understanding your encrypted messages",
			}
			s.Send(msg)
		} else {
			h.State.Encryption[m.Name] = secret
			time.Sleep(1 * time.Second) // wait for textual to get situated
			encryptedMessage, err := dh1080.Enc([]byte("I've stared an encrypted chat, just message .encoff to remove our group key"), []byte(secret))
			if err != nil {
				return
			}
			encryptedBaseMessage := base64.StdEncoding.EncodeToString(encryptedMessage)
			infoMessage := &irc.Message{
				Command:  irc.PRIVMSG,
				Params:   []string{m.Name},
				Trailing: "+OK *" + encryptedBaseMessage,
			}
			s.Send(infoMessage)
		}
	}
}
