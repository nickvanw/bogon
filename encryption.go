package bogon

import (
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nickvanw/bogon/dh1080"
	"github.com/nickvanw/ircx"
	irc "gopkg.in/sorcix/irc.v1"
)

const encOff string = ".encoff"

// NewEncryption creates a new in-memory conversation encryption storage
func NewEncryption() Encryption {
	return &MemoryEncryption{sessions: map[string]encSession{}}
}

// Encryption is an interface used to create and store encrypted chat sessions
type Encryption interface {
	New(string, string)
	Rename(string, string)
	Check(string) (string, bool)
	Stop(string)
}

// MemoryEncryption is an in-memory store of ongoing encrypted conversations
type MemoryEncryption struct {
	sync.Mutex
	sessions map[string]encSession
}

type encSession struct {
	Started time.Time
	Key     string
}

// New creates a new session with the specified user and key
func (s *MemoryEncryption) New(user, key string) {
	s.Lock()
	s.sessions[user] = encSession{Key: key, Started: time.Now()}
	s.Unlock()
}

// Rename transfers the session to a different user
func (s *MemoryEncryption) Rename(oldname, newname string) {
	s.Lock()
	if d, ok := s.sessions[oldname]; ok {
		delete(s.sessions, oldname)
		s.sessions[newname] = d
	}
	s.Unlock()
}

// Check returns the encryption key, as well as a boolean if there is a
// current encrypted session
func (s *MemoryEncryption) Check(user string) (string, bool) {
	if key, ok := s.sessions[user]; ok {
		return key.Key, true
	}
	return "", false
}

// Stop removes an encrypted session with the specified user
func (s *MemoryEncryption) Stop(user string) {
	s.Lock()
	delete(s.sessions, user)
	s.Unlock()
}

func (c *Client) encryptionMsgHandler(s ircx.Sender, m *irc.Message) {
	if dm, ok := c.decodedMessage(m); ok && dm.Trailing == ".encoff" {
		c.state.Encryption().Stop(m.Name)
		out := "I've removed the key for this conversation"
		msg := &irc.Message{
			Command:  irc.PRIVMSG,
			Params:   []string{m.Name},
			Trailing: out,
		}
		s.Send(msg)
	}
}

func (c *Client) encryptionStartHandler(s ircx.Sender, m *irc.Message) {
	data := strings.Split(m.Trailing, " ")
	if len(data) < 2 {
		return
	}
	if data[0] == "DH1080_INIT" {
		enc := dh1080.New()
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
			c.state.Encryption().New(m.Name, secret)
			time.Sleep(1 * time.Second) // wait for their client to get situated
			out := fmt.Sprintf("I've stared an encrypted chat, just message %s to remove our group key", encOff)
			encryptedMessage, err := dh1080.Enc([]byte(out), []byte(secret))
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
