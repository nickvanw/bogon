package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

func init() {
	AddPlugin("Reddit", "(?i)^\\.r(eddit)?u(ser)?$", MessageHandler(Reddit), false, false)
}

const layout = "January 2, 2006"

func Reddit(msg *Message) {
	if len(msg.Params) < 2 {
		msg.Return("Give me a user!")
		return
	}
	user := msg.Params[1]
	url := fmt.Sprintf("http://www.reddit.com/user/%s/about.json", user)
	data, _ := getSite(url)
	var response RedditUser
	json.Unmarshal(data, &response)
	if len(response.Data.Name) < 1 {
		msg.Return("Couldn't find that person!")
		return
	}
	//name := response.Data.Name
	created := response.Data.CreatedUtc
	createdInt := int64(created)
	ParsedTime := time.Unix(createdInt, 0)
	createdString := ParsedTime.Format(layout)
	outString := fmt.Sprintf("%s with %v comment and %v link karma, created on %s", response.Data.Name, humanize.Comma(int64(response.Data.CommentKarma)), humanize.Comma(int64(response.Data.LinkKarma)), createdString)
	msg.Return(outString)

}

type RedditUser struct {
	Data struct {
		CommentKarma     float64     `json:"comment_karma"`
		Created          float64     `json:"created"`
		CreatedUtc       float64     `json:"created_utc"`
		HasMail          interface{} `json:"has_mail"`
		HasModMail       interface{} `json:"has_mod_mail"`
		HasVerifiedEmail bool        `json:"has_verified_email"`
		ID               string      `json:"id"`
		IsFriend         bool        `json:"is_friend"`
		IsGold           bool        `json:"is_gold"`
		IsMod            bool        `json:"is_mod"`
		LinkKarma        float64     `json:"link_karma"`
		Name             string      `json:"name"`
		Over18           bool        `json:"over_18"`
	} `json:"data"`
	Kind string `json:"kind"`
}
