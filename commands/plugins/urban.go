package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/util"
)

const urbanapi = "http://api.urbandictionary.com/v0/define?page=1&term="

var urbanCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.urban$")
	return urbanTitle, out, urbanLookup, defaultOptions
}

func urbanLookup(msg commands.Message, ret commands.MessageFunc) string {
	search := strings.Join(msg.Params[1:], " ")
	url := fmt.Sprintf("%s%s", urbanapi, util.URLEncode(search))
	data, err := util.Fetch(url)
	if err != nil {
		return "Error trying to find that!"
	}
	var ud urbanInfo
	if err := json.Unmarshal(data, &ud); err != nil {
		return "Urban Dictionary gave me a bad response!"
	}
	if len(ud.List) == 0 {
		return "I didn't get anything for that!"
	}

	return fmt.Sprintf("%s: %s", util.Bold(search), ud.List[0].Definition)

}

type urbanInfo struct {
	List []struct {
		Author      string  `json:"author"`
		CurrentVote string  `json:"current_vote"`
		Defid       float64 `json:"defid"`
		Definition  string  `json:"definition"`
		Example     string  `json:"example"`
		Permalink   string  `json:"permalink"`
		ThumbsDown  float64 `json:"thumbs_down"`
		ThumbsUp    float64 `json:"thumbs_up"`
		Word        string  `json:"word"`
	} `json:"list"`
	ResultType string        `json:"result_type"`
	Sounds     []interface{} `json:"sounds"`
	Tags       []string      `json:"tags"`
}
