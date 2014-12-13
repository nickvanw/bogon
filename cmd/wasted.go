package cmd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

func init() {
	AddPlugin("Wasted", "(?i)^\\.wasted$", MessageHandler(Wasted), false, false)
}

const wastedurl = "http://reddit.com/r/wastedgifs.json"

func Wasted(msg *Message) {
	data, err := getSite(wastedurl)

	if err != nil {
		msg.Return("Error getting reddit data!")
		return
	}

	var rp RedditPage
	json.Unmarshal(data, &rp)

	rand.Seed(time.Now().UTC().UnixNano())
	randIndex := rand.Intn(len(rp.Data.Children))

	wasteddata := rp.Data.Children[randIndex].ChildData
	var out string

	wasteddata.Url = fixLink(wasteddata.Url)

	if wasteddata.Nsfw {
		out = fmt.Sprintf("(NSFW) %s: %s", wasteddata.Title, wasteddata.Url)
	} else {
		out = fmt.Sprintf("%s: %s", wasteddata.Title, wasteddata.Url)
	}

	msg.Return(out)
}

// Converts a link to a direct link if it's on imgur
func fixLink(link string) string {
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
