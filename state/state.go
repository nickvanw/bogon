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
	bot.AddCallback(irc.JOIN, ircx.Callback{Handler: &JoinHandler{Bot: bot, State: s}})
	bot.AddCallback(irc.PART, ircx.Callback{Handler: &PartHandler{Bot: bot, State: s}})
	bot.AddCallback(irc.QUIT, ircx.Callback{Handler: &QuitHandler{Bot: bot, State: s}})
	bot.AddCallback(irc.KICK, ircx.Callback{Handler: &KickHandler{Bot: bot, State: s}})
	bot.AddCallback(irc.MODE, ircx.Callback{Handler: &ModeHandler{Bot: bot, State: s}})
	bot.AddCallback(irc.RPL_TOPIC, ircx.Callback{Handler: &TopicHandler{Bot: bot, State: s}})
	bot.AddCallback(irc.RPL_NAMREPLY, ircx.Callback{Handler: &NamesHandler{Bot: bot, State: s}})
	bot.AddCallback(irc.ERR_NICKNAMEINUSE, ircx.Callback{Handler: &NickTakenHandler{Bot: bot, State: s}})
	bot.AddCallback(irc.NICK, ircx.Callback{Handler: &NickHandler{Bot: bot, State: s}})
	bot.AddCallback(irc.INVITE, ircx.Callback{Handler: &InviteHandler{Bot: bot, State: s}})
	bot.AddCallback(irc.NOTICE, ircx.Callback{Handler: &EncryptionHandler{Bot: bot, State: s}})
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
