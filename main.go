package main

import (
	"flag"
	"log"
	"os"

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
	user     = flag.String("user", "ircx", "User to send to IRC server")
	password = flag.String("password", "", "Password to send to irc server")
	channels = flag.String("chan", "#test", "Channels to join")
	redis    = flag.String("redis", "127.0.0.1:6379", "Redis Port")
)

func main() {
	flag.Parse()
	*password = os.Getenv("PASS")
	newServer := Bot{Bot: ircx.WithLogin(*server, *name, *user, *password), State: &state.State{Encryption: map[string]string{}}}
	if err := newServer.Connect(); err != nil {
		log.Panicln("Unable to dial IRC Server ", err)
	}
	newServer.State.Name = *name
	RegisterCoreHandlers(newServer.Bot, newServer.State)
	state.RegisterStateHandlers(newServer.Bot, newServer.State)
	newServer.CallbackLoop()
	log.Println("Exiting..")
}
