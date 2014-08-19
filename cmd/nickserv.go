package cmd

import (
	"fmt"
	"strings"

	"github.com/nickvanw/bogon/state"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

type NickServHandler struct {
	Bot   *ircx.Bot
	State *state.State
}

func (h *NickServHandler) Handle(s irc.Sender, m *irc.Message) {
	username, u_ok := GetConfig("NickservUser")
	password, p_ok := GetConfig("NickservPass")
	nickserv, ns_ok := GetConfig("Nickserv")
	phrase, ph_ok := GetConfig("NickservPhrase")
	if u_ok && p_ok && ns_ok && ph_ok && m.Name == nickserv && strings.Contains(m.Trailing, phrase) {
		msg := &irc.Message{
			Command:  irc.PRIVMSG,
			Params:   []string{"NickServ"},
			Trailing: fmt.Sprintf("IDENTIFY %s %s", username, password),
		}
		s.Send(msg)
	}
}
