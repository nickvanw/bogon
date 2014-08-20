package cmd

import (
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
)

func init() {
	AddPlugin("Bookmark", "(?i)^\\.b(ook)?m(ark)?$", MessageHandler(Bookmark), false, false)
	AddPlugin("RemoveBM", "(?i)^\\.delbm$", MessageHandler(RemoveBM), false, false)
}

func Bookmark(msg *Message) {
	if len(msg.Params) < 3 {
		msg.Return("Usage: .bm [key] [message to bookmark]")
	} else {
		key := msg.Params[1]
		if len(key) < 2 {
			msg.Return("Sorry, I couldn't insert that.")
			return
		}
		message := strings.Join(msg.Params[2:], " ")
		success := InsertBookmark(key, message)
		if success {
			msg.Return("Successfully added")
		} else {
			msg.Return("Sorry, I couldn't insert that.")
		}
	}
}

func RemoveBM(msg *Message) {
	if len(msg.Params) < 2 {
		msg.Return("Usage: .delbm [key]")
	} else {
		key := msg.Params[1]
		success := RemoveBookmark(key)
		if success {
			msg.Return("Successfully Removed")
		} else {
			msg.Return("I don't think that exists, or I couldn't remove it, sorry!")
		}
	}
}

func InsertBookmark(key string, msg string) bool {
	key = strings.ToLower(key)
	conn := pool.Get()
	defer conn.Close()
	if isCmd(fmt.Sprintf(".%s", key)) {
		return false
	} else {
		bmhash := fmt.Sprintf("%s:%s", PREFIX, "bm")
		_, err := conn.Do("hset", bmhash, key, msg)
		if err != nil {
			return false
		}
		return true
	}
}

func isCmd(key string) bool {
	for _, v := range Commands {
		if v.Command.MatchString(key) && !v.Raw {
			return true
		}
	}
	return false
}

func RemoveBookmark(key string) bool {
	key = strings.ToLower(key)
	conn := pool.Get()
	defer conn.Close()
	bmhash := fmt.Sprintf("%s:%s", PREFIX, "bm")
	data, err := redis.Int(conn.Do("hdel", bmhash, key))
	if err != nil || data == 0 {
		return false
	}
	return true
}
