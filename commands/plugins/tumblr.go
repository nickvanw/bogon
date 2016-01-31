package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/bogon/commands/util"
)

var titpCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.titp$")
	return titpTitle, out, titpTumblr, defaultOptions
}

var tumblrCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.tumblr$")
	return tumblrTitle, out, tumblrLookup, defaultOptions
}

func titpTumblr(msg commands.Message, ret commands.MessageFunc) string {
	msg.Params = []string{msg.Params[0], "thisisthinprivilege.tumblr.com"}
	return tumblrLookup(msg, ret)
}

func tumblrLookup(msg commands.Message, ret commands.MessageFunc) string {
	if len(msg.Params) < 2 {
		return "I need .tumblr [blog-url]"
	}
	blog := msg.Params[1]
	key, ok := config.Get("TUMBLR_API")
	if !ok {
		return ""
	}

	req := fmt.Sprintf("http://api.tumblr.com/v2/blog/%s/posts/text?api_key=%s", blog, key)
	data, err := util.Fetch(req)
	if err != nil {
		return "Unable to get data from tumblr"
	}
	var resp tumblrData
	if err := json.Unmarshal(data, &resp); err != nil {
		return "Tumblr gave me a bad response!"
	}
	if resp.Meta.Status != 200 {
		return "Tumblr returned an Error"
	}
	if len(resp.Response.Posts) < 1 {
		return "I didn't find any posts!"
	}
	post := resp.Response.Posts[0]
	return fmt.Sprintf("%s: %v @ %s", post.BlogName, post.Title, post.ShortURL)
}

type tumblrData struct {
	Meta struct {
		Msg    string `json:"msg"`
		Status int64  `json:"status"`
	} `json:"meta"`
	Response struct {
		Blog struct {
			Ask          bool   `json:"ask"`
			AskAnon      bool   `json:"ask_anon"`
			AskPageTitle string `json:"ask_page_title"`
			Description  string `json:"description"`
			IsNsfw       bool   `json:"is_nsfw"`
			Name         string `json:"name"`
			Posts        int64  `json:"posts"`
			ShareLikes   bool   `json:"share_likes"`
			Title        string `json:"title"`
			Updated      int64  `json:"updated"`
			URL          string `json:"url"`
		} `json:"blog"`
		Posts []struct {
			BlogName     string        `json:"blog_name"`
			Body         string        `json:"body"`
			Date         string        `json:"date"`
			Format       string        `json:"format"`
			Highlighted  []interface{} `json:"highlighted"`
			ID           int64         `json:"id"`
			IsSubmission bool          `json:"is_submission"`
			NoteCount    int64         `json:"note_count"`
			PostAuthor   interface{}   `json:"post_author"`
			PostURL      string        `json:"post_url"`
			ReblogKey    string        `json:"reblog_key"`
			ShortURL     string        `json:"short_url"`
			Slug         string        `json:"slug"`
			SourceTitle  string        `json:"source_title"`
			SourceURL    string        `json:"source_url"`
			State        string        `json:"state"`
			Tags         []interface{} `json:"tags"`
			Timestamp    int64         `json:"timestamp"`
			Title        interface{}   `json:"title"`
			Type         string        `json:"type"`
		} `json:"posts"`
		TotalPosts int64 `json:"total_posts"`
	} `json:"response"`
}
