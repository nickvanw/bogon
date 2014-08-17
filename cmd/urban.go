package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
)

func init() {
	AddPlugin("Urban", "(?i)^\\.urban$", MessageHandler(Urban), false, false)
}

const urbanapi = "http://api.urbandictionary.com/v0/define?page=1&term="

func Urban(msg *Message) {
	search := strings.Join(msg.Params[1:], " ")
	url := fmt.Sprintf("%s%s", urbanapi, urlencode(search))
	data, err := getSite(url)
	if err != nil {
		msg.Return("Error trying to find that!")
		return
	}
	var ud UrbanInfo
	json.Unmarshal(data, &ud)
	if len(ud.List) == 0 {
		msg.Return("I didn't get anything for that!")
		return
	} else {
		out_fact := ud.List[0]
		out := fmt.Sprintf("%s: %s", bold(search), out_fact.Definition)
		msg.Return(out)
	}
}

type UrbanInfo struct {
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
