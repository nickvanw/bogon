package main

import (
	"flag"
	"log"
	"math"
	"os"
	"time"

	"github.com/nickvanw/bogon/cmd"
	"github.com/nickvanw/bogon/state"
	"github.com/nickvanw/ircx"
)

type Bot struct {
	*ircx.Bot
	*state.State
}

var (
	name        = flag.String("name", "ircx", "Nick to use in IRC")
	server      = flag.String("server", "chat.freenode.org:6667", "Host:Port to connect to")
	user        = flag.String("user", "ircx", "User to send to IRC server")
	password    = flag.String("password", "", "Password to send to irc server")
	channels    = flag.String("chan", "#test", "Channels to join")
	redisServer = flag.String("redis", "127.0.0.1:6379", "Redis Host:Port")
	config      = flag.String("config", "config.toml", "Config file")
)

func main() {
	flag.Parse()
	cmd.InitCommand(*config, *redisServer)
	*password = os.Getenv("PASS")
	newServer := Bot{Bot: ircx.WithLogin(*server, *name, *user, *password), State: &state.State{Encryption: map[string]string{}}}
	tries := float64(1)
	for err := newServer.Connect(); err != nil; err = newServer.Connect() {
		duration := time.Duration(math.Pow(2.0, tries)*200) * time.Millisecond
		log.Println("Unable to connect to", *server, "- waiting", duration)
		time.Sleep(duration)
		tries++
	}
	newServer.State.Name = *name
	RegisterCoreHandlers(newServer.Bot, newServer.State)
	state.RegisterStateHandlers(newServer.Bot, newServer.State)
	newServer.CallbackLoop()
	log.Println("Exiting..")
}
