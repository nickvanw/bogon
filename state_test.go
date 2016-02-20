package bogon

import (
	"reflect"
	"sync"
	"testing"

	"github.com/nickvanw/ircx"
	"github.com/sorcix/irc"
)

const name = "test_bot"

var (
	emptyChan     = &Channel{Modes: make(map[rune]struct{}), Users: map[string]*User{}}
	defaultPrefix = &irc.Prefix{Name: name, User: "~ircx", Host: "184.152.8.252"}
	serverPrefix  = &irc.Prefix{Host: "someserver.biz"}
)

func TestBasicState(t *testing.T) {
	var tt = []struct {
		name      string
		data      []*irc.Message
		wantState func() *MemoryState
		wantOut   []*irc.Message
	}{
		{
			name: "joining a channel creates a channel in state",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string(nil), Trailing: "#test", EmptyTrailing: false},
			},
			wantState: func() *MemoryState {
				s := NewState(name, WithRejoin())
				s.Channels = map[string]*Channel{"#test": emptyChan}
				return s
			},
		},
		{
			name: "inviting the bot results in it sending a join",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "INVITE", Params: []string{"ircx", "#test2"}, Trailing: "", EmptyTrailing: false},
			},
			wantState: func() *MemoryState {
				return NewState(name, WithRejoin())
			},
			wantOut: []*irc.Message{&irc.Message{Command: "JOIN", Params: []string{"#test2"}, Trailing: "", EmptyTrailing: false}},
		},
		{
			name: "joining and then kicking results in the channel being deleted from state; join being sent",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string{"#test2"}, Trailing: "", EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "KICK", Params: []string{"#test2", name}},
			},
			wantState: func() *MemoryState {
				return NewState(name, WithRejoin())
			},
			wantOut: []*irc.Message{&irc.Message{Command: "JOIN", Params: []string{"#test2"}, Trailing: "", EmptyTrailing: false}},
		},
		{
			name: "recieving a NICK changes my name",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "NICK", Params: []string{}, Trailing: "butt", EmptyTrailing: false},
			},
			wantState: func() *MemoryState {
				s := NewState("butt", WithRejoin())
				return s
			},
		},
		{
			name: "joining a channel creates a channel, sending NAMES populates the channel names",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string(nil), Trailing: "#test3", EmptyTrailing: false},
				&irc.Message{Prefix: serverPrefix, Command: "353", Params: []string{name, "=", "#test3"}, Trailing: "butt @Guest53 Guest101 " + name, EmptyTrailing: false},
			},
			wantState: func() *MemoryState {
				s := NewState(name, WithRejoin())
				ch := &Channel{Modes: make(map[rune]struct{})}
				ch.Users = map[string]*User{
					"butt":     &User{Modes: map[rune]struct{}{}},
					"guest53":  &User{Modes: map[rune]struct{}{'o': struct{}{}}},
					"guest101": &User{Modes: map[rune]struct{}{}},
					name:       &User{Modes: map[rune]struct{}{}},
				}
				s.Channels = map[string]*Channel{"#test3": ch}
				return s
			},
		},
		{
			name: "joining a channel creates a channel, sending NAMES, user sending NICK changes their name",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string(nil), Trailing: "#test3", EmptyTrailing: false},
				&irc.Message{Prefix: serverPrefix, Command: "353", Params: []string{name, "=", "#test3"}, Trailing: "butt @Guest53 Guest101 " + name, EmptyTrailing: false},
				&irc.Message{Prefix: &irc.Prefix{Name: "guest53", User: "~ircx", Host: "8.8.8.8"}, Command: "NICK", Params: []string{}, Trailing: "guest54", EmptyTrailing: false},
			},
			wantState: func() *MemoryState {
				s := NewState(name, WithRejoin())
				ch := &Channel{Modes: make(map[rune]struct{})}
				ch.Users = map[string]*User{
					"butt":     &User{Modes: map[rune]struct{}{}},
					"guest54":  &User{Modes: map[rune]struct{}{'o': struct{}{}}},
					"guest101": &User{Modes: map[rune]struct{}{}},
					name:       &User{Modes: map[rune]struct{}{}},
				}
				s.Channels = map[string]*Channel{"#test3": ch}
				return s
			},
		},
		{
			name: "joining a channel creates a channel, sending TOPIC populates the topic",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string(nil), Trailing: "#test3", EmptyTrailing: false},
				&irc.Message{Prefix: serverPrefix, Command: "332", Params: []string{name, "#test3"}, Trailing: "1 2 3 4 5", EmptyTrailing: false},
			},
			wantState: func() *MemoryState {
				s := NewState(name, WithRejoin())
				ch := &Channel{Modes: make(map[rune]struct{}), Users: make(map[string]*User)}
				ch.Topic = "1 2 3 4 5"
				s.Channels = map[string]*Channel{"#test3": ch}
				return s
			},
		},
		{
			name: "joining a channel creates a channel, sending NAMES populates the channel names, KICK removes one of them",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string(nil), Trailing: "#test3", EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "353", Params: []string{name, "=", "#test3"}, Trailing: "butt @Guest53 Guest101 " + name, EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "KICK", Params: []string{"#test3", "Guest101"}, EmptyTrailing: true},
			},
			wantState: func() *MemoryState {
				s := NewState(name, WithRejoin())
				ch := &Channel{Modes: make(map[rune]struct{})}
				ch.Users = map[string]*User{
					"butt":    &User{Modes: map[rune]struct{}{}},
					"guest53": &User{Modes: map[rune]struct{}{'o': struct{}{}}},
					name:      &User{Modes: map[rune]struct{}{}},
				}
				s.Channels = map[string]*Channel{"#test3": ch}
				return s
			},
		},
		{
			name: "joining a channel creates a channel, sending NAMES populates the channel names, KICK removes one of them, JOIN adds a new one",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string(nil), Trailing: "#test3", EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "353", Params: []string{name, "=", "#test3"}, Trailing: "butt @Guest53 Guest101 " + name, EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "KICK", Params: []string{"#test3", "Guest101"}, EmptyTrailing: true},
				&irc.Message{Prefix: &irc.Prefix{Name: "guest104", User: "~ircx", Host: "8.8.8.8"}, Command: "JOIN", Params: []string{"#test3"}, EmptyTrailing: true},
			},
			wantState: func() *MemoryState {
				s := NewState(name, WithRejoin())
				ch := &Channel{Modes: make(map[rune]struct{})}
				ch.Users = map[string]*User{
					"butt":     &User{Modes: map[rune]struct{}{}},
					"guest104": &User{Modes: map[rune]struct{}{}},
					"guest53":  &User{Modes: map[rune]struct{}{'o': struct{}{}}},
					name:       &User{Modes: map[rune]struct{}{}},
				}
				s.Channels = map[string]*Channel{"#test3": ch}
				return s
			},
		},
		{
			name: "joining a channel creates a channel, sending NAMES populates the channel names, QUIT removes one of them",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string(nil), Trailing: "#test3", EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "353", Params: []string{name, "=", "#test3"}, Trailing: "butt @Guest53 Guest101 " + name, EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string(nil), Trailing: "#test4", EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "353", Params: []string{name, "=", "#test4"}, Trailing: "Guest53 Guest101 " + name, EmptyTrailing: false},
				&irc.Message{Prefix: &irc.Prefix{Name: "guest101", User: "~ircx", Host: "8.8.8.8"}, Command: "QUIT", Params: []string(nil), EmptyTrailing: true},
			},
			wantState: func() *MemoryState {
				s := NewState(name, WithRejoin())
				ch1 := &Channel{Modes: make(map[rune]struct{})}
				ch1.Users = map[string]*User{
					"butt":    &User{Modes: map[rune]struct{}{}},
					"guest53": &User{Modes: map[rune]struct{}{'o': struct{}{}}},
					name:      &User{Modes: map[rune]struct{}{}},
				}
				ch2 := &Channel{Modes: make(map[rune]struct{})}
				ch2.Users = map[string]*User{
					"guest53": &User{Modes: map[rune]struct{}{}},
					name:      &User{Modes: map[rune]struct{}{}},
				}

				s.Channels = map[string]*Channel{"#test3": ch1, "#test4": ch2}
				return s
			},
		},
		{
			name: "joining a channel creates a channel, sending NAMES populates the channel names, PART removes one of them",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string(nil), Trailing: "#test3", EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "353", Params: []string{name, "=", "#test3"}, Trailing: "butt @Guest53 Guest101 " + name, EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string(nil), Trailing: "#test4", EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "353", Params: []string{name, "=", "#test4"}, Trailing: "Guest53 Guest101 " + name, EmptyTrailing: false},
				&irc.Message{Prefix: &irc.Prefix{Name: "guest101", User: "~ircx", Host: "8.8.8.8"}, Command: "PART", Params: []string{"#test4"}, EmptyTrailing: true},
			},
			wantState: func() *MemoryState {
				s := NewState(name, WithRejoin())
				ch1 := &Channel{Modes: make(map[rune]struct{})}
				ch1.Users = map[string]*User{
					"butt":     &User{Modes: map[rune]struct{}{}},
					"guest101": &User{Modes: map[rune]struct{}{}},
					"guest53":  &User{Modes: map[rune]struct{}{'o': struct{}{}}},
					name:       &User{Modes: map[rune]struct{}{}},
				}
				ch2 := &Channel{Modes: make(map[rune]struct{})}
				ch2.Users = map[string]*User{
					"guest53": &User{Modes: map[rune]struct{}{}},
					name:      &User{Modes: map[rune]struct{}{}},
				}

				s.Channels = map[string]*Channel{"#test3": ch1, "#test4": ch2}
				return s
			},
		},
		{
			name: "joining a channel creates a channel, sending NAMES populates the channel names, MODE changes ops/channel",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string(nil), Trailing: "#test3", EmptyTrailing: false},
				&irc.Message{Prefix: serverPrefix, Command: "353", Params: []string{name, "=", "#test3"}, Trailing: "butt @Guest53 Guest101 " + name, EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "MODE", Params: []string{"#test3", "+oo", "butt", "guest101"}, Trailing: "", EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "MODE", Params: []string{"#test3", "-o", "guest53"}, Trailing: "", EmptyTrailing: false},
				&irc.Message{Prefix: defaultPrefix, Command: "MODE", Params: []string{"#test3", "+m"}, Trailing: "", EmptyTrailing: false},
			},
			wantState: func() *MemoryState {
				s := NewState(name, WithRejoin())
				ch := &Channel{Modes: map[rune]struct{}{'m': struct{}{}}}
				ch.Users = map[string]*User{
					"butt":     &User{Modes: map[rune]struct{}{'o': struct{}{}}},
					"guest53":  &User{Modes: map[rune]struct{}{}},
					"guest101": &User{Modes: map[rune]struct{}{'o': struct{}{}}},
					name:       &User{Modes: map[rune]struct{}{}},
				}
				s.Channels = map[string]*Channel{"#test3": ch}
				return s
			},
		},
		{
			name: "joining and then parting results in the channel being deleted from state",
			data: []*irc.Message{
				&irc.Message{Prefix: defaultPrefix, Command: "JOIN", Params: []string{"#test2"}, EmptyTrailing: true},
				&irc.Message{Prefix: defaultPrefix, Command: "PART", Params: []string{"#test2", name}, EmptyTrailing: true},
			},
			wantState: func() *MemoryState {
				return NewState(name, WithRejoin())
			},
		},
		{
			name: "the server telling me my name is already taken results in me appending a |",
			data: []*irc.Message{
				&irc.Message{Prefix: serverPrefix, Command: irc.ERR_NICKNAMEINUSE, Params: []string{name}, EmptyTrailing: true},
			},
			wantState: func() *MemoryState {
				return NewState(name+"|", WithRejoin())
			},
			wantOut: []*irc.Message{&irc.Message{Command: irc.NICK, Params: []string{name + "|"}, EmptyTrailing: false}},
		},
	}

	for _, v := range tt {
		t.Logf("Starting test: %s", v.name)
		b, c := bogonFactory()
		for _, v := range v.data {
			c.Recv <- v
		}
		close(c.Recv)
		<-c.Done
		wantState := v.wantState()
		if !reflect.DeepEqual(b.state, wantState) {
			t.Logf("want: %#v\n", wantState)
			t.Logf(" got: %#v\n", b.state)
			t.Fatalf("state was not equal")
		}
		if !reflect.DeepEqual(v.wantOut, c.Out) {
			t.Logf("want: %#v\n", v.wantOut)
			t.Logf(" got: %#v\n", c.Out)
			t.Fatalf("returned data was not equal")
		}
	}
}

// bogonFactory begins a bogon listening on a fake socket that allows
// the passing of arbitrary irc messages in to test for how state works.
func bogonFactory() (*Client, *c) {
	cn := &c{Recv: make(chan *irc.Message), Done: make(chan struct{}), handlers: map[string][]ircx.Handler{}}
	bot, err := New(cn, name, []string{})
	if err != nil {
		panic(err)
	}
	bot.Start()
	return bot, cn
}

type c struct {
	sync.Mutex
	Recv     chan *irc.Message
	Done     chan struct{}
	Out      []*irc.Message
	handlers map[string][]ircx.Handler
}

func (f *c) Connect() error {
	return nil
}

func (f *c) HandleLoop() {
	go func() {
		for {
			select {
			case m, ok := <-f.Recv:
				if !ok {
					close(f.Done)
					return
				}
				handlers, ok := f.handlers[m.Command]
				if !ok {
					return
				}
				for _, h := range handlers {
					h.Handle(f, m)
				}
			}
		}
	}()
}

// Send sends the specified message
func (f *c) Send(msg *irc.Message) error {
	f.Lock()
	f.Out = append(f.Out, msg)
	f.Unlock()
	return nil
}

func (f *c) HandleFunc(t string, fn func(ircx.Sender, *irc.Message)) {
	f.handlers[t] = append(f.handlers[t], ircx.HandlerFunc(fn))
}
