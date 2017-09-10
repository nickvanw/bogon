package reddit

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/dustin/go-humanize"
)

var reddit = new(redditAuth)

func userLookup(user string) string {
	url := fmt.Sprintf(userUrl, user)
	data, err := reddit.fetch(url)
	if err != nil {
		return "Could not fetch data from reddit!"
	}
	var ret redditUser
	if err := json.Unmarshal(data, &ret); err != nil {
		return "Could not read data from reddit!"
	}
	return fmt.Sprintf("%s with %v comment and %v link karma, created on %s", ret.Data.Name,
		humanize.Comma(int64(ret.Data.CommentKarma)), humanize.Comma(int64(ret.Data.LinkKarma)),
		time.Unix(int64(ret.Data.CreatedUtc), 0).Format("January 2, 2006"))
}

type redditUser struct {
	Data struct {
		CommentKarma float64 `json:"comment_karma"`
		CreatedUtc   float64 `json:"created_utc"`
		LinkKarma    float64 `json:"link_karma"`
		Name         string  `json:"name"`
	} `json:"data"`
}

func subLookup(sub string) string {
	url := fmt.Sprintf(subUrl, sub)
	data, err := reddit.fetch(url)
	if err != nil {
		return "Could not fetch data from reddit!"
	}
	var ret redditSub
	if err := json.Unmarshal(data, &ret); err != nil {
		return "Could not read data from reddit!"
	}

	postcount := len(ret.Data.Children)
	if postcount == 0 {
		return "This subreddit is empty!"
	}

	rand.Seed(time.Now().UTC().UnixNano())
	randIndex := rand.Intn(postcount)

	post := ret.Data.Children[randIndex].ChildData
	post.URL = fixImgurLink(post.URL)

	if post.Nsfw {
		return fmt.Sprintf("(NSFW) %s: %s", post.Title, post.URL)
	}
	return fmt.Sprintf("%s: %s", post.Title, post.URL)
}

type redditSub struct {
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

// Converts a link to a direct link if it's on imgur
func fixImgurLink(link string) string {
	originalurl, err := url.Parse(link)

	if err != nil || originalurl.Host != "imgur.com" {
		return link
	}

	return fmt.Sprintf("http://i.imgur.com%s.gif", originalurl.Path)
}
