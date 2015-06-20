package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func init() {
	AddPlugin("Wiki", "(?i)^\\.wiki(pedia)?$", MessageHandler(Wiki), false, false)
}

const wikiurl = "http://en.wikipedia.org/w/api.php?format=json"

func Wiki(msg *Message) {
	if len(msg.Params) < 2 {
		return
	}

	query := url.QueryEscape(strings.Join(msg.Params[1:], " "))
	url := fmt.Sprintf("%s&action=query&prop=extracts|info&exintro=&explaintext=&inprop=url&indexpageids=&redirects=&titles=%s", wikiurl, query)

	data, err := getSite(url)

	if err != nil {
		msg.Return("Error contacting Wikipedia API!")
		return
	}

	var resp WikiResponse
	json.Unmarshal(data, &resp)

	if len(resp.Query.PageIds) < 1 {
		msg.Return("Page not found!")
		return
	}

	id := resp.Query.PageIds[0]
	page := resp.Query.Pages[id]

	content := page.Extract
	if len(content) > 350 {
		content = fmt.Sprintf("%s...", content[:350])
	}

	out := fmt.Sprintf("%s %s", page.Url, content)
	msg.Return(out)
}

type WikiResponse struct {
	Query struct {
		PageIds []string        `json:"pageids"`
		Pages   map[string]Page `json:"pages"`
	} `json:"query"`
}

type Page struct {
	PageId  int    `json:"pageid"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
	Url     string `json:"fullurl"`
}
