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

	// Set up command system (redis, config file)
	cmd.InitCommand(*config, *redisServer)

	// check if we're using a pasword to connect to IRC
	*password = os.Getenv("PASS")

	// set up the Bot we'll be connecting to IRC
	newBot := ircx.WithLogin(*server, *name, *user, *password)

	// create the bot's state and initialize it
	newState := &state.State{Encryption: map[string]string{}, Name: *name, Channels: map[string]*state.Channel{}}
	newState.InitState()

	// create our combination of bot and state
	newServer := Bot{Bot: newBot, State: newState}

	// set the try count to 1
	tries := float64(1)

	// loop forever with exponential backup to connect
	for err := newServer.Connect(); err != nil; err = newServer.Connect() {
		duration := time.Duration(math.Pow(2.0, tries)*200) * time.Millisecond
		log.Println("Unable to connect to", *server, "- waiting", duration)
		time.Sleep(duration)
		tries++
	}

	// Register the essential handlers
	RegisterCoreHandlers(newServer.Bot, newState)

	// register the state handlers
	newState.RegisterStateHandlers(newServer.Bot)

	// start processing callbacks
	newServer.HandleLoop()
}
