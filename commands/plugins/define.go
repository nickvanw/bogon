package plugins

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/bogon/commands/util"
)

var defineCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.def(ine|inition)?$")
	return defineTitle, out, defineWord, defaultOptions
}

func defineWord(msg commands.Message, ret commands.MessageFunc) string {
	word := url.QueryEscape(strings.Join(msg.Params[1:], " "))
	apiKey, avail := config.Get("WORDNIK_API")
	if avail != true {
		return ""
	}
	url := fmt.Sprintf("http://api.wordnik.com/v4/word.json/%s/definitions?includeRelated=false&api_key=%s&includeTags=false&limit=1&useCanonical=true", word, apiKey)
	data, err := util.Fetch(url)
	if err != nil {
		return "Unable to look up the definition of that word"
	}
	var response []defineReturn
	if err := json.Unmarshal(data, &response); err != nil {
		return "Wordnik gave me a bad answer for that word"
	}
	if len(response) > 0 {
		return fmt.Sprintf("%s (%s): %s", response[0].Word, response[0].PartOfSpeech, response[0].Text)
	}
	return "I couldn't find that word, sorry!"
}

type defineReturn struct {
	Word         string
	Text         string
	PartOfSpeech string
}
