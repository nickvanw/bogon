package bogon

import (
	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

// Client contains a single IRC bot
type Client struct {
	bot   bot
	state state

	commands []command
}

type bot interface {
	Connect() error
	HandleLoop()
	HandleFunc(string, func(s ircx.Sender, m *irc.Message))
}

// State models all of the interactions an IRC server can have with the state
// of the client, which is a lot.
// Any implementation of State should be careful to avoid races
type state interface {
	NewUser(channel string, user string)
	RemoveUser(channel string, name string)
	QuitUser(name string)
	RenameUser(oldname, newname string)

	NewChannel(name string)
	GetChan(channel string) (*Channel, error)
	RemoveChannel(name string)

	Name() string
	SetName(newname string)

	ParseModes(modes []string)
	Encryption() Encryption
	Rejoin() bool
}

// New accepts an underlying ircx connection and a list of channels to join
// when connected.
func New(bot bot, name string, channels []string) (*Client, error) {
	client := &Client{bot: bot}
	client.state = NewState(name, WithRejoin())
	client.registerStateHandlers()
	client.registerCoreHandlers(channels)
	return client, nil
}

// Connect begins an IRC connection
func (c *Client) Connect() error {
	return c.bot.Connect()
}

// Start begins the blocking callback loop
func (c *Client) Start() {
	c.bot.HandleLoop()
}

// AdminSocket starts the admin socket at the specified file in a goroutine
// if a file exists there, it will automatically be removed
func (c *Client) AdminSocket(fn string) {
	go c.adminSocket(fn)
}

func (c *Client) registerCoreHandlers(channels []string) {
	c.bot.HandleFunc(irc.RPL_WELCOME, registerConnect(channels))
	c.bot.HandleFunc(irc.PING, pingHandler)
}

func pingHandler(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command:  irc.PONG,
		Params:   m.Params,
		Trailing: m.Trailing,
	})
}

func registerConnect(c []string) func(s ircx.Sender, m *irc.Message) {
	return func(s ircx.Sender, m *irc.Message) {
		s.Send(&irc.Message{
			Command: irc.JOIN,
			Params:  c,
		})
	}
}
