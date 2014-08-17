package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/fzzy/radix/redis"
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
	redis, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err != nil || isCmd(fmt.Sprintf(".%s", key)) {
		return false
	} else {
		bmhash := fmt.Sprintf("%s:%s", PREFIX, "bm")
		ret := redis.Cmd("hset", bmhash, key, msg)
		if ret.Err != nil {
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
	redis, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err != nil {
		return false
	} else {
		bmhash := fmt.Sprintf("%s:%s", PREFIX, "bm")
		ret := redis.Cmd("hdel", bmhash, key)
		val, reterr := ret.Int()
		if reterr != nil || val == 0 {
			return false
		}
		return true
	}
}
