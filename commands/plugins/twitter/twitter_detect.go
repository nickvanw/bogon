package twitter

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/util"
)

const twitterDetectTitle = "twitter_detect"

func (a *API) rawHandler(msg commands.Message, ret commands.MessageFunc) string {
	for _, v := range msg.Params {
		u, err := url.Parse(v)
		if err != nil {
			continue
		}
		// basically, only this: https://twitter.com/molly_knight/status/970808912818577410
		if (u.Scheme == "https" || u.Scheme == "http") && (u.Host == "twitter.com") {
			p, tweet := path.Split(u.Path)
			if strings.HasSuffix(p, "/status/") {
				tweet, err := strconv.Atoi(tweet)
				if err != nil {
					return ""
				}
				data, err := a.api.GetTweet(int64(tweet), url.Values{})
				if err != nil {
					return ""
				}
				var time string
				when, err := data.CreatedAtTime()
				if err != nil {
					time = "Unknown"
				} else {
					time = when.Local().Format("Mon Jan 2, 2006 @ 3:04pm")
				}
				return fmt.Sprintf("%s: %s [%s rt, %s fav] on %s", data.User.ScreenName, util.StripNewLines(data.Text),
					humanize.Comma(int64(data.RetweetCount)), humanize.Comma(int64(data.FavoriteCount)), time)
			}
		}
	}
	return ""
}

// RawTwitterHandler registers a handler to look for tweets in a list of strings
func (a *API) RawTwitterHandler() commands.RegisterFunc {
	return func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
		return twitterDetectTitle, nil, a.rawHandler, commands.Options{Raw: true}
	}
}
