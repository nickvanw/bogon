package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func init() {
	AddPlugin("Define", "(?i)^\\.def(ine|inition)?$", MessageHandler(Define), false, false)
}

type DefineReturn struct {
	Word         string
	Text         string
	PartOfSpeech string
}

func Define(msg *Message) {
	word := url.QueryEscape(strings.Join(msg.Params[1:], " "))
	WORDNIK, avail := GetConfig("Define")
	if avail != true {
		fmt.Println("Wordnik API Key not available!")
		return
	}
	url := fmt.Sprintf("http://api.wordnik.com/v4/word.json/%s/definitions?includeRelated=false&api_key=%s&includeTags=false&limit=1&useCanonical=true", word, WORDNIK)
	data, err := getSite(url)
	if err != nil {
		msg.Return("Error!")
		return
	}
	var response []DefineReturn
	json.Unmarshal(data, &response)
	if len(response) > 0 {
		out := fmt.Sprintf("%s (%s): %s", response[0].Word, response[0].PartOfSpeech, response[0].Text)
		msg.Return(out)
	} else {
		msg.Return("I couldn't find that word, sorry!")
	}
}
