package cmd

import (
	"fmt"
	"regexp"

	"github.com/BurntSushi/toml"
)

var (
	Commands   []Cmd
	ConfigData map[string]string
)

func init() {
	_, err := toml.DecodeFile("config.toml", &ConfigData)
	if err != nil {
		fmt.Printf("Unable to load config file, this means multiple commands will not work: %s\r\n", err)
		return
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

func GetConfig(conf string) (string, bool) {
	conf, ok := ConfigData[conf]
	return conf, ok
}
