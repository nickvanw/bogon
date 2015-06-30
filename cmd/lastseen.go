package cmd

import (
	"bytes"
	"fmt"
	"time"

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

	diff := time.Now().Unix() - seen

	years := diff / 60 / 60 / 24 / 365
	diff -= years * 60 * 60 * 24 * 365

	days := diff / 60 / 60 / 24
	diff -= days * 60 * 60 * 24

	hours := diff / 60 / 60
	diff -= hours * 60 * 60

	minutes := diff / 60
	diff -= minutes * 60

	var buffer bytes.Buffer

	if years > 1 {
		buffer.WriteString(fmt.Sprintf("%d years, ", years))
	} else if years == 1 {
		buffer.WriteString(fmt.Sprintf("%d year, ", years))
	}

	if days > 1 {
		buffer.WriteString(fmt.Sprintf("%d days, ", days))
	} else if days == 1 {
		buffer.WriteString(fmt.Sprintf("%d day, ", days))
	}

	if hours > 1 {
		buffer.WriteString(fmt.Sprintf("%d hours, ", hours))
	} else if hours == 1 {
		buffer.WriteString(fmt.Sprintf("%d hour, ", hours))
	}

	if minutes != 1 {
		buffer.WriteString(fmt.Sprintf("%d minutes ago", minutes))
	} else {
		buffer.WriteString(fmt.Sprintf("%d minute ago", minutes))
	}

	return buffer.String()
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
