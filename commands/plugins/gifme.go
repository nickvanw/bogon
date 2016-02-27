package plugins

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/gifgo"
)

var gifmeCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.gifme$")
	return gifmeTitle, out, gifMe, defaultOptions
}

func gifMe(msg commands.Message, ret commands.MessageFunc) string {
	var opts []gifgo.OptFunc
	if key, ok := config.Get("GIPHY_KEY"); ok {
		opts = append(opts, gifgo.APIKey(key))
	}
	client, err := gifgo.New(opts...)
	if err != nil {
		return "Unable to talk to Giphy"
	}
	var randQuery gifgo.RandomReq
	if len(msg.Params) > 1 {
		randQuery.Tag = strings.Join(msg.Params[1:], " ")
	}
	gif, err := client.Random(randQuery)
	if err != nil {
		fmt.Println(err)
		return "Giphy returned an error"
	}
	desc := "(No Description)"
	if gif.Data.Caption != "" {
		desc = gif.Data.Caption
	}
	return fmt.Sprintf("%s: %s - Powered by GIPHY", desc, gif.Data.ImageURL)
}
