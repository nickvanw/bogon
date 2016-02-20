package bogon

import (
	"errors"
	"log"
	"strings"
	"sync"
)

var (
	// ErrChannelNotFound is returned when the specified channel is not being state tracked
	ErrChannelNotFound = errors.New("channel not found")
	// ErrUserNotFound is returned when the specified user cannot be found
	ErrUserNotFound = errors.New("user not found")
)

type stateConf func(s *MemoryState)

// WithRejoin automatically rejoins channels when kicked
func WithRejoin() stateConf {
	return func(s *MemoryState) {
		s.rejoin = true
	}
}

// NewState returns an in-memory state for a single IRC connection
func NewState(name string, conf ...stateConf) *MemoryState {
	s := &MemoryState{
		name:       name,
		encryption: NewEncryption(),
		Channels:   map[string]*Channel{},
	}
	for _, v := range conf {
		v(s)
	}
	return s
}

// MemoryState tracks the state of an IRC connection in memory
type MemoryState struct {
	sync.RWMutex
	Channels   map[string]*Channel
	Motd       string
	encryption Encryption
	Admin      string
	name       string
	rejoin     bool
}

// Channel represents a single IRC channel the client is currently in
type Channel struct {
	sync.RWMutex
	Topic string
	Users map[string]*User
	Modes map[rune]struct{}
}

// User represents a single user in a single channel
// if the same user overlaps in multiple channels with the connection,
// there will be multiple User structs in each channel
type User struct {
	sync.RWMutex
	Modes map[rune]struct{}
}

// Encryption returns the underlying encryption tracking
func (s *MemoryState) Encryption() Encryption {
	return s.encryption
}

// Name returns the current name of the IRC connection
func (s *MemoryState) Name() string {
	s.RLock()
	defer s.RUnlock()
	return s.name
}

// SetName updates the name of the connected client
func (s *MemoryState) SetName(n string) {
	s.Lock()
	s.name = n
	s.Unlock()
}

// Rejoin returns true if the connection should auto-rejoin when kicked
func (s *MemoryState) Rejoin() bool {
	s.RLock()
	defer s.RUnlock()
	return s.rejoin
}

// GetChan returns the channel of specified string
func (s *MemoryState) GetChan(channel string) (*Channel, error) {
	s.RLock()
	defer s.RUnlock()
	if findChannel, ok := s.Channels[strings.ToLower(channel)]; ok {
		return findChannel, nil
	}
	return nil, ErrChannelNotFound
}

// QuitUser will remove the User from every channel the state is tracking
func (s *MemoryState) QuitUser(name string) {
	s.Lock()
	for _, channel := range s.Channels {
		channel.RemoveUser(name)
	}
	s.Unlock()
}

// RemoveChannel removes the specified channel from state tracking
func (s *MemoryState) RemoveChannel(name string) {
	s.Lock()
	delete(s.Channels, strings.ToLower(name))
	s.Unlock()
}

// RemoveUser removes the specified user from the specified channel
// if they are presen
func (s *MemoryState) RemoveUser(channel, name string) {
	remChannel, err := s.GetChan(channel)
	if err != nil {
		return
	}
	remChannel.RemoveUser(name)
}

// NewUser registers the user in the specified channel
func (s *MemoryState) NewUser(channel, user string) {
	addChannel, err := s.GetChan(channel)
	if err != nil {
		return
	}
	addChannel.NewUser(user)
}

// NewChannel creates a channel of the specified name, adding it to the state
func (s *MemoryState) NewChannel(name string) {
	s.Lock()
	s.Channels[strings.ToLower(name)] = &Channel{Modes: make(map[rune]struct{}), Users: map[string]*User{}}
	s.Unlock()
}

// RenameUser renames the user from old to new name in every channel the state is tracking
func (s *MemoryState) RenameUser(oldname, newname string) {
	s.Lock()
	for _, channel := range s.Channels {
		channel.Lock()
		if val, ok := channel.Users[strings.ToLower(oldname)]; ok {
			channel.Users[strings.ToLower(newname)] = val
			delete(channel.Users, strings.ToLower(oldname))
		}
		channel.Unlock()
	}
	s.Encryption().Rename(oldname, newname)
	s.Unlock()
}

// ParseModes returns the modes that are prepended to a user
// such as o/v
func (s *MemoryState) ParseModes(modes []string) {
	channel, err := s.GetChan(modes[0])
	if err != nil {
		// todo(nick): These are modes for myself
		log.Println("I got modes for a channel I am not in")
		return
	}
	modeString := modes[1]
	modeArgs := modes[2:]
	var plus bool
	for _, v := range modeString {
		switch v {
		case '+':
			plus = true
		case '-':
			plus = false
		case 'q', 'a', 'o', 'h', 'v':
			nick := modeArgs[0]
			user, err := channel.GetUser(nick)
			if err != nil {
				continue
			}
			user.Lock()
			if plus {
				user.Modes[v] = struct{}{}
			} else {
				delete(user.Modes, v)
			}
			user.Unlock()
			modeArgs = modeArgs[1:]
		default:
			channel.Lock()
			if plus {
				channel.Modes[v] = struct{}{}
			} else {
				delete(channel.Modes, v)
			}
			channel.Unlock()
		}
	}
}

// SetTopic sets the topic of the specified channel
func (c *Channel) SetTopic(topic string) {
	c.Lock()
	c.Topic = topic
	c.Unlock()
}

// GetUser gets the specified user from the channel
func (c *Channel) GetUser(name string) (*User, error) {
	c.RLock()
	defer c.RUnlock()
	if findUser, ok := c.Users[strings.ToLower(name)]; ok {
		return findUser, nil
	}
	return nil, ErrUserNotFound
}

// RemoveUser removes the user from the channel if they exist
func (c *Channel) RemoveUser(name string) {
	c.Lock()
	delete(c.Users, strings.ToLower(name))
	c.Unlock()
}

// NewUser adds one-or-more users to the specified channel
func (c *Channel) NewUser(users ...string) {
	for _, user := range users {
		c.RemoveUser(user)
		modes := make(map[rune]struct{})
		switch rune(user[0]) {
		case '~', '&', '@', '%', '+':
			modes[symbolToRune[string(user[0])]] = struct{}{}
			user = user[1:]
		}
		c.Lock()
		c.Users[strings.ToLower(user)] = &User{Modes: modes}
		c.Unlock()
	}
}
