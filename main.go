package main

import (
	"flag"
	"log"

	"github.com/nickvanw/bogon/state"
	"github.com/nickvanw/ircx"
)

type Bot struct {
	*ircx.Bot
	*state.State
}

var (
	name     = flag.String("name", "ircx", "Nick to use in IRC")
	server   = flag.String("server", "chat.freenode.org:6667", "Host:Port to connect to")
	channels = flag.String("chan", "#test", "Channels to join")
	redis    = flag.String("redis", "127.0.0.1:6379", "Redis Port")
)

func main() {
	flag.Parse()
	newServer := Bot{Bot: ircx.Classic(*server, *name), State: &state.State{}}
	if err := newServer.Connect(); err != nil {
		log.Panicln("Unable to dial IRC Server ", err)
	}
	newServer.State.Name = *name
	RegisterCoreHandlers(newServer.Bot, newServer.State)
	state.RegisterStateHandlers(newServer.Bot, newServer.State)
	newServer.CallbackLoop()
	log.Println("Exiting..")
}
