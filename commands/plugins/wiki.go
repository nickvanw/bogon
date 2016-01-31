package plugins

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/util"
)

const wikiurl = "http://en.wikipedia.org/w/api.php?format=json"

var wikiCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.wiki(pedia)?$")
	return wikiTitle, out, wikiLookup, defaultOptions
}

func wikiLookup(msg commands.Message, ret commands.MessageFunc) string {
	if len(msg.Params) < 2 {
		return ""
	}

	query := url.QueryEscape(strings.Join(msg.Params[1:], " "))
	url := fmt.Sprintf("%s&action=query&prop=extracts|info&exintro=&explaintext=&inprop=url&indexpageids=&redirects=&titles=%s", wikiurl, query)

	data, err := util.Fetch(url)

	if err != nil {
		return "Error contacting Wikipedia API!"
	}

	var resp wikiResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "Wikipedia gave me a bad answer!"
	}

	if len(resp.Query.PageIDs) < 1 {
		return "Wikipedia page not found!"
	}

	id := resp.Query.PageIDs[0]
	page := resp.Query.Pages[id]

	content := page.Extract
	if len(content) > 350 {
		content = fmt.Sprintf("%s...", content[:350])
	}

	return fmt.Sprintf("%s %s", page.URL, content)
}

type wikiResponse struct {
	Query struct {
		PageIDs []string            `json:"pageids"`
		Pages   map[string]wikiPage `json:"pages"`
	} `json:"query"`
}

type wikiPage struct {
	PageID  int    `json:"pageid"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
	URL     string `json:"fullurl"`
}
