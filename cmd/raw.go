package cmd

import (
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
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
	go HandleLastUrl(msg)
	go HandleLastSeen(msg)
}

func HandleRedis(msg *Message) {
	bm, ex := GetIfExistsBookmark(msg.Params[0])
	if ex {
		msg.Return(bm)
	}
}

func GetIfExistsBookmark(begin string) (string, bool) {
	conn := pool.Get()
	defer conn.Close()
	if len(begin) < 1 {
		return "", false
	}
	if strings.Contains(CMD, string(begin[0])) {
		begin = strings.ToLower(begin[1:])
	} else {
		return "", false
	}
	testString := fmt.Sprintf("%s:%s", PREFIX, "bm")
	ret, err := redis.String(conn.Do("hget", testString, begin))
	if err != nil {
		return "", false
	} else {
		dispMsg := fmt.Sprintf("[%s]: %s", begin, ret)
		return dispMsg, true
	}
}

func GetBookmark(id string) string {
	conn := pool.Get()
	defer conn.Close()
	testString := fmt.Sprintf("%s:%s", PREFIX, "bm")
	data, err := redis.String(conn.Do("hget", testString, strings.ToLower(id[1:])))
	if err != nil {
		return "false"
	} else {
		return data
	}
}
