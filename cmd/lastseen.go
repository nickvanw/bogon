package cmd

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/garyburd/redigo/redis"
)

func init() {
	AddPlugin("LastSeen", "(?i)^\\.(last|seen)$", MessageHandler(LastSeen), false, false)
}

func LastSeen(msg *Message) {
	if len(msg.Params) < 2 {
		msg.Return("Usage: .last [nick]")
	} else {
		nick := msg.Params[1]
		seen := GetLastSeen(nick)

		msg.Return(fmt.Sprintf("%s was last seen %s", nick, seen))
	}
}

func GetLastSeen(nick string) string {
	conn := pool.Get()
	defer conn.Close()

	cmd := fmt.Sprintf("%s:%s", PREFIX, "lastseen")
	seen, err := redis.Int64(conn.Do("hget", cmd, nick))

	if err != nil {
		return "never"
	}

	return humanize.Time(time.Unix(seen, 0))
}

func HandleLastSeen(msg *Message) {
	nick := msg.User.Name
	now := time.Now().Unix()

	conn := pool.Get()
	defer conn.Close()

	cmd := fmt.Sprintf("%s:%s", PREFIX, "lastseen")
	_, err := conn.Do("hset", cmd, nick, now)

	if err != nil {
		fmt.Printf("Could not save last seen timestamp for %s\n", nick)
	}
}
