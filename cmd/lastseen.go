package cmd

import (
	"fmt"
	"strings"
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
		seen := GetLastSeen(msg)
		msg.Return(fmt.Sprintf("%s was last seen %s", msg.Params[1], seen))
	}
}

func GetLastSeen(msg *Message) string {
	nick := msg.Params[1]
	channel := msg.To
	key := fmt.Sprintf("%s%s", nick, channel)
	msgkey := fmt.Sprintf("%s:msg", key)

	conn := pool.Get()
	defer conn.Close()
	cmd := fmt.Sprintf("%s:%s", PREFIX, "lastseen")

	seen, err := redis.Int64(conn.Do("hget", cmd, key))
	if err != nil {
		return "never"
	}
	seenstr := humanize.Time(time.Unix(seen, 0))

	lastmsg, err := redis.String(conn.Do("hget", cmd, msgkey))
	if err != nil {
		return seenstr
	}

	return fmt.Sprintf("%s, %q", seenstr, lastmsg)
}

func HandleLastSeen(msg *Message) {
	nick := msg.User.Name
	channel := msg.To
	key := fmt.Sprintf("%s%s", nick, channel)
	msgkey := fmt.Sprintf("%s:msg", key)

	now := time.Now().Unix()
	lastmsg := strings.Join(msg.Params, " ")

	conn := pool.Get()
	defer conn.Close()
	cmd := fmt.Sprintf("%s:%s", PREFIX, "lastseen")

	_, err := conn.Do("hset", cmd, key, now)
	if err != nil {
		fmt.Printf("Could not save last seen timestamp for %s in %s\n", nick, channel)
		return
	}

	_, err = conn.Do("hset", cmd, msgkey, lastmsg)
	if err != nil {
		fmt.Printf("Could not save last message for %s in %s\n", nick, channel)
	}
}
