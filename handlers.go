package main

import (
	"github.com/nickvanw/bogon/cmd"
	"github.com/nickvanw/bogon/state"
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

func RegisterCoreHandlers(bot *ircx.Bot, state *state.State) {
	// Add the initial join and ping handlers
	bot.HandleFunc(irc.RPL_WELCOME, RegisterConnect)
	bot.HandleFunc(irc.PING, PingHandler)

	// Add the command handler for channel commands
	bot.Handle(irc.PRIVMSG, &cmd.CommandHandler{Bot: bot, State: state})
	bot.Handle(irc.NOTICE, &cmd.NickServHandler{Bot: bot, State: state})
}

func PingHandler(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   m.Params,
		Trailing: m.Trailing,
	})
}

func RegisterConnect(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{*channels},
	})
}
