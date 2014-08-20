package cmd

import (
	"log"
	"regexp"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/garyburd/redigo/redis"
)

var (
	Commands   []Cmd
	ConfigData map[string]string
)

var (
	pool     *redis.Pool
	myConfig string
)

func InitCommand(config, redis string) {
	myConfig = config
	pool = newPool(redis)
	if err := LoadConfig(config); err != nil {
		log.Println("Unable to get config file, some commands and features will not work")
	}
}

func AddPlugin(name string, regex string, function MessageHandler, raw bool, admin bool) {
	reg := regexp.MustCompile(regex)
	cmd := Cmd{
		Name:         name,
		Command:      reg,
		Function:     function,
		Raw:          raw,
		RequireAdmin: admin,
	}
	Commands = append(Commands, cmd)
}

type Cmd struct {
	Name         string
	Command      *regexp.Regexp
	Function     MessageHandler
	RequireAdmin bool
	Raw          bool
}

func LoadConfig(file string) error {
	_, err := toml.DecodeFile(file, &ConfigData)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig(conf string) (string, bool) {
	if ConfigData == nil {
		return "", false
	}
	conf, ok := ConfigData[conf]
	return conf, ok
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
