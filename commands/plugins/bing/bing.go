package bing

import (
	"regexp"
	"strings"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
)

const (
	searchTitle = "bing"
	imageTitle  = "image_search"
)

// BingSearch registers a handler to search bing search
var BingSearch = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.(google|bing)$")
	return searchTitle, out, bingSearch, commands.Options{}
}

// ImageSearch registers a handler to search bing images and upload the first result to S3
var ImageSearch = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.im(a)?g(e)?$")
	return imageTitle, out, bingImageSearch, commands.Options{}
}

func bingSearch(msg commands.Message, ret commands.MessageFunc) string {
	return bingProcess(msg, ret, bingSearchProcess{})
}
func bingImageSearch(msg commands.Message, ret commands.MessageFunc) string {
	return bingProcess(msg, ret, bingImageProcess{})
}

func bingProcess(msg commands.Message, ret commands.MessageFunc, p bingProcesser) string {
	token, ok := config.Get("BING_API")
	if !ok {
		return ""
	}
	msg.Params = append(msg.Params, "nsfw")
	query := strings.Join(msg.Params[1:], " ")
	out, err := bingAPIFetch(query, token, p)
	if err != nil {
		return "Unable to execute that search, sorry."
	}
	return out
}
