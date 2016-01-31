package plugins

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"time"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/util"
)

const (
	subredditURL = "http://reddit.com/r/%s.json"
)

var subredditCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.(sub)?r(eddit)?$")
	return subredditTitle, out, subredditLookup, defaultOptions
}

func subredditLookup(msg commands.Message, ret commands.MessageFunc) string {
	if len(msg.Params) < 2 {
		return "Give me a subreddit!"
	}

	sub := msg.Params[1]
	url := fmt.Sprintf(subredditURL, sub)

	data, err := util.Fetch(url)

	if err != nil {
		return "Error getting reddit data!"
	}

	var rp redditPage
	json.Unmarshal(data, &rp)

	postcount := len(rp.Data.Children)
	if postcount == 0 {
		return "This subreddit looks empty!"
	}

	rand.Seed(time.Now().UTC().UnixNano())
	randIndex := rand.Intn(postcount)

	redditdata := rp.Data.Children[randIndex].ChildData
	redditdata.URL = fixImgurLink(redditdata.URL)

	if redditdata.Nsfw {
		return fmt.Sprintf("(NSFW) %s: %s", redditdata.Title, redditdata.URL)
	}
	return fmt.Sprintf("%s: %s", redditdata.Title, redditdata.URL)
}

// Converts a link to a direct link if it's on imgur
func fixImgurLink(link string) string {
	originalurl, err := url.Parse(link)

	if err != nil || originalurl.Host != "imgur.com" {
		return link
	}

	return fmt.Sprintf("http://i.imgur.com%s.gif", originalurl.Path)
}

type redditPage struct {
	Data struct {
		Children []struct {
			ChildData struct {
				Title string `json:"title"`
				Nsfw  bool   `json:"over_18"`
				URL   string `json:"url"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}
