package state

import (
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

var SymbolToRune = map[string]rune{
	"@": 'o',
	"+": 'v',
	"%": 'h',
	"&": 'a',
	"~": 'q',
}

func (s *State) RegisterStateHandlers(bot *ircx.Bot) {
	bot.Handle(irc.JOIN, &JoinHandler{Bot: bot, State: s})
	bot.Handle(irc.PART, &PartHandler{Bot: bot, State: s})
	bot.Handle(irc.QUIT, &QuitHandler{Bot: bot, State: s})
	bot.Handle(irc.KICK, &KickHandler{Bot: bot, State: s, Rejoin: true})
	bot.Handle(irc.MODE, &ModeHandler{Bot: bot, State: s})
	bot.Handle(irc.RPL_TOPIC, &TopicHandler{Bot: bot, State: s})
	bot.Handle(irc.RPL_NAMREPLY, &NamesHandler{Bot: bot, State: s})
	bot.Handle(irc.ERR_NICKNAMEINUSE, &NickTakenHandler{Bot: bot, State: s})
	bot.Handle(irc.NICK, &NickHandler{Bot: bot, State: s})
	bot.Handle(irc.INVITE, &InviteHandler{Bot: bot, State: s})
	bot.Handle(irc.NOTICE, &EncryptionHandler{Bot: bot, State: s})
}

type State struct {
	sync.Mutex
	Channels   map[string]*Channel
	Motd       string
	Name       string
	Encryption map[string]string
	Password   string
	Admin      string
}

func (s *State) InitState() {
	s.Password = os.Getenv("ADMIN_PASS")
}

type Channel struct {
	Topic   string
	Users   map[string]*User
	Modes   map[rune]struct{}
	LastUrl string
}

type User struct {
	Modes map[rune]struct{}
}

func (s *State) GetChan(channel string) (*Channel, error) {
	if findChannel, ok := s.Channels[strings.ToLower(channel)]; ok {
		return findChannel, nil
	}
	return nil, errors.New("Not a channel!")
}

func (c *Channel) GetUser(name string) (*User, error) {
	if findUser, ok := c.Users[strings.ToLower(name)]; ok {
		return findUser, nil
	} else {
		return nil, errors.New("Not a user!")
	}
}

func (s *State) RemoveChannel(name string) {
	s.Lock()
	defer s.Unlock()
	delete(s.Channels, strings.ToLower(name))
}

func (s *State) RemoveUser(channel string, name string) {
	s.Lock()
	defer s.Unlock()
	remChannel, err := s.GetChan(channel)
	if err != nil {
		return
	}
	delete(remChannel.Users, strings.ToLower(name))
}

func (c *Channel) RemoveUser(name string) {
	delete(c.Users, strings.ToLower(name))
}

func (s *State) NewUser(channel string, user string) {
	s.Lock()
	defer s.Unlock()
	addChannel, err := s.GetChan(channel)
	if err != nil {
		return
	}
	addChannel.NewUser(user)
}

func (c *Channel) NewUser(user string) {
	c.RemoveUser(user)
	modes := make(map[rune]struct{})
	switch rune(user[0]) {
	case '~', '&', '@', '%', '+':
		modes[SymbolToRune[string(user[0])]] = struct{}{}
		user = user[1:]
	}
	c.Users[strings.ToLower(user)] = &User{Modes: modes}
}

func (s *State) NewChannel(name string) {
	s.Lock()
	defer s.Unlock()
	s.Channels[strings.ToLower(name)] = &Channel{Modes: make(map[rune]struct{}), Users: map[string]*User{}}
}

func (c *Channel) Ops() int {
	num := 0
	for _, v := range c.Users {
		if _, ok := v.Modes['o']; ok {
			num++
		}
	}
	return num
}

func (c *Channel) Voice() int {
	num := 0
	for _, v := range c.Users {
		if _, ok := v.Modes['v']; ok {
			num++
		}
	}
	return num
}
