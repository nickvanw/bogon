package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/fzzy/radix/redis"
)

func init() {
	AddPlugin("Raw", "", MessageHandler(Raw), true, false)
}

const PREFIX = "ManaGo"
const CMD = "."

func Raw(msg *Message) {
	go HandleRedis(msg)
	go HandleYoutube(msg.Params, msg)
	go HandleSpotify(msg)
}

func HandleRedis(msg *Message) {
	redis, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err != nil {
		fmt.Println("[ERR]: Redis Error!")
	} else {
		bm, ex := GetIfExistsBookmark(msg.Params[0], redis)
		if ex {
			msg.Return(bm)
		}
	}
}

func GetIfExistsBookmark(begin string, rc *redis.Client) (string, bool) {
	if len(begin) < 1 {
		return "", false
	}
	if strings.Contains(CMD, string(begin[0])) {
		begin = strings.ToLower(begin[1:])
	} else {
		return "", false
	}
	testString := fmt.Sprintf("%s:%s", PREFIX, "bm")
	ret, err := rc.Cmd("hget", testString, begin).Str()
	if err != nil {
		return "", false
	} else {
		dispMsg := fmt.Sprintf("[%s]: %s", begin, ret)
		return dispMsg, true
	}
}

func GetBookmark(id string, rc *redis.Client) string {
	testString := fmt.Sprintf("%s:%s", PREFIX, "bm")
	data, err := rc.Cmd("hget", testString, strings.ToLower(id[1:])).Str()
	if err != nil {
		return "false"
	} else {
		return data
	}
}
