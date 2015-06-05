package cmd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

func init() {
	AddPlugin("SubReddit", "(?i)^\\.(sub)?r(eddit)?$", MessageHandler(SubReddit), false, false)
}

const baseurl = "http://reddit.com/r/%s.json"

func SubReddit(msg *Message) {
	if len(msg.Params) < 2 {
		msg.Return("Give me a subreddit!")
		return
	}

	sub := msg.Params[1]
	url := fmt.Sprintf(baseurl, sub);

	data, err := getSite(url)

	if err != nil {
		msg.Return("Error getting reddit data!")
		return
	}

	var rp RedditPage
	json.Unmarshal(data, &rp)

	rand.Seed(time.Now().UTC().UnixNano())
	randIndex := rand.Intn(len(rp.Data.Children))

	redditdata := rp.Data.Children[randIndex].ChildData
	var out string

	redditdata.Url = fixImgurLink(redditdata.Url)

	if redditdata.Nsfw {
		out = fmt.Sprintf("(NSFW) %s: %s", redditdata.Title, redditdata.Url)
	} else {
		out = fmt.Sprintf("%s: %s", redditdata.Title, redditdata.Url)
	}

	msg.Return(out)
}

// Converts a link to a direct link if it's on imgur
func fixImgurLink(link string) string {
	originalurl, err := url.Parse(link)

	if err != nil || originalurl.Host != "imgur.com" {
		return link
	}

	return fmt.Sprintf("http://i.imgur.com%s.gif", originalurl.Path)
}

type RedditPage struct {
	Data struct {
		Children []struct {
			ChildData struct {
				Title string `json:"title"`
				Nsfw  bool   `json:"over_18"`
				Url   string `json:"url"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}
