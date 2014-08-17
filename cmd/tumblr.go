package cmd

import (
	"encoding/json"
	"fmt"
	"log"
)

func init() {
	AddPlugin("Tumblr", "(?i)^\\.tumblr$", MessageHandler(Tumblr), false, false)
	AddPlugin("TiTP", "(?i)^\\.titp$", MessageHandler(TiTP), false, false)
}

func TiTP(msg *Message) {
	msg.Params = []string{msg.Params[0], "thisisthinprivilege.tumblr.com"}
	Tumblr(msg)
}

func Tumblr(msg *Message) {
	if len(msg.Params) < 2 {
		msg.Return("I need .tumblr [blog-url]")
		return
	}
	blog := msg.Params[1]
	key, ok := GetConfig("Tumblr")
	if !ok {
		log.Println("I don't have a Tumblr API key")
		return
	}

	req := fmt.Sprintf("http://api.tumblr.com/v2/blog/%s/posts/text?api_key=%s", blog, key)
	data, err := getSite(req)
	if err != nil {
		msg.Return("Unable to get data from tumblr")
		return
	}
	var resp TumblrData
	json.Unmarshal(data, &resp)
	if resp.Meta.Status != 200 {
		msg.Return("Tumblr returned an Error")
		return
	}
	if len(resp.Response.Posts) < 1 {
		msg.Return("I didn't find any posts!")
		return
	}
	post := resp.Response.Posts[0]
	output := fmt.Sprintf("%s: %s @ %s", post.BlogName, post.Title, post.ShortUrl)
	msg.Return(output)
}

type TumblrData struct {
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
			Url          string `json:"url"`
		} `json:"blog"`
		Posts []struct {
			BlogName     string        `json:"blog_name"`
			Body         string        `json:"body"`
			Date         string        `json:"date"`
			Format       string        `json:"format"`
			Highlighted  []interface{} `json:"highlighted"`
			Id           int64         `json:"id"`
			IsSubmission bool          `json:"is_submission"`
			NoteCount    int64         `json:"note_count"`
			PostAuthor   interface{}   `json:"post_author"`
			PostUrl      string        `json:"post_url"`
			ReblogKey    string        `json:"reblog_key"`
			ShortUrl     string        `json:"short_url"`
			Slug         string        `json:"slug"`
			SourceTitle  string        `json:"source_title"`
			SourceUrl    string        `json:"source_url"`
			State        string        `json:"state"`
			Tags         []interface{} `json:"tags"`
			Timestamp    int64         `json:"timestamp"`
			Title        interface{}   `json:"title"`
			Type         string        `json:"type"`
		} `json:"posts"`
		TotalPosts int64 `json:"total_posts"`
	} `json:"response"`
}
