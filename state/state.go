package state

import (
	"errors"
	"log"
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

func RegisterStateHandlers(bot *ircx.Bot, state *State) {
	bot.AddCallback(irc.JOIN, ircx.Callback{Handler: &JoinHandler{Bot: bot, State: state}})
	bot.AddCallback(irc.PART, ircx.Callback{Handler: &PartHandler{Bot: bot, State: state}})
	bot.AddCallback(irc.QUIT, ircx.Callback{Handler: &QuitHandler{Bot: bot, State: state}})
	bot.AddCallback(irc.KICK, ircx.Callback{Handler: &KickHandler{Bot: bot, State: state}})
	bot.AddCallback(irc.MODE, ircx.Callback{Handler: &ModeHandler{Bot: bot, State: state}})
	bot.AddCallback(irc.RPL_TOPIC, ircx.Callback{Handler: &TopicHandler{Bot: bot, State: state}})
	bot.AddCallback(irc.RPL_NAMREPLY, ircx.Callback{Handler: &NamesHandler{Bot: bot, State: state}})
}

type State struct {
	sync.Mutex
	Channels []*Channel
	Motd     string
	Name     string
}

type Channel struct {
	Name  string
	Topic string
	Users []*User
	Modes map[rune]struct{}
}

type User struct {
	Name  string
	Modes map[rune]struct{}
}

func (s *State) GetChan(channel string) (*Channel, error) {
	for _, chan_try := range s.Channels {
		if strings.ToLower(chan_try.Name) == strings.ToLower(channel) {
			return chan_try, nil
		}
	}
	return &Channel{}, errors.New("Not a channel!")
}

func (c *Channel) GetUser(name string) *User {
	for _, user := range c.Users {
		if strings.ToLower(user.Name) == strings.ToLower(name) {
			return user
		}
	}
	return &User{}
}

func (s *State) RemoveChannel(name string) {
	s.Lock()
	defer s.Unlock()
	for num, pchan := range s.Channels {
		if strings.ToLower(pchan.Name) == strings.ToLower(name) {
			s.Channels = append(s.Channels[:num], s.Channels[num+1:]...)
		}
	}
}

func (s *State) RemoveUser(channel string, name string) {
	s.Lock()
	defer s.Unlock()
	remChannel, err := s.GetChan(channel)
	if err != nil {
		log.Println("I tried to remove a user form a channel I have no record of")
		return
	}
	remChannel.RemoveUser(name)
}

func (c *Channel) RemoveUser(name string) {
	for num, user := range c.Users {
		if strings.ToLower(user.Name) == strings.ToLower(name) {
			c.Users = append(c.Users[:num], c.Users[num+1:]...)
		}
	}
}

func (s *State) NewUser(channel string, user string) {
	s.Lock()
	defer s.Unlock()
	addChannel, err := s.GetChan(channel)
	if err != nil {
		log.Println("User joined a channel I have no record of")
		return
	}
	addChannel.NewUser(user)
}

func (c *Channel) NewUser(user string) {
	modes := make(map[rune]struct{})
	switch rune(user[0]) {
	case '~', '&', '@', '%', '+':
		modes[SymbolToRune[string(user[0])]] = struct{}{}
		user = user[1:]
	}
	c.Users = append(c.Users, &User{Name: user, Modes: modes})
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
