package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/util"
)

const (
	redditLayout = "January 2, 2006"
)

var redditCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.r(eddit)?u(ser)?$")
	return redditTitle, out, redditLookup, defaultOptions
}

func redditLookup(msg commands.Message, ret commands.MessageFunc) string {
	if len(msg.Params) < 2 {
		return "Usage: .ru [username]"
	}
	user := msg.Params[1]
	url := fmt.Sprintf("http://www.reddit.com/user/%s/about.json", user)
	data, err := util.Fetch(url)
	if err != nil {
		return "Unable to fetch the users' information from reddit"
	}
	var response redditUser
	if err := json.Unmarshal(data, &response); err != nil {
		return "Reddit gave me a bad response"
	}
	if len(response.Data.Name) < 1 {
		return "Couldn't find that person!"
	}
	//name := response.Data.Name
	created := response.Data.CreatedUtc
	createdInt := int64(created)
	ParsedTime := time.Unix(createdInt, 0)
	createdString := ParsedTime.Format(redditLayout)
	return fmt.Sprintf("%s with %v comment and %v link karma, created on %s",
		response.Data.Name, humanize.Comma(int64(response.Data.CommentKarma)),
		humanize.Comma(int64(response.Data.LinkKarma)), createdString)

}

type redditUser struct {
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
