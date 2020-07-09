package twitter

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/ChimeraCoder/anaconda"
	"github.com/dustin/go-humanize"
	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/bogon/commands/util"
)

var (
	// ErrMissingTokens is returned when the environment does not have the necessary
	// API tokens
	ErrMissingTokens = errors.New("one or more api tokens is missing")
)

const twitterTitle = "twitter"

// API is a twitter API client with a command handler
type API struct {
	api *anaconda.TwitterApi
}

// NewFromEnv checks config for the necessary tokens, returning an error
// if there are any missing
func NewFromEnv() (*API, error) {
	a := new(API)
	consumerKey, cKeyOk := config.Get("TWITTER_CONSUMER_KEY")
	consumerSecret, cSecOk := config.Get("TWITTER_CONSUMER_SECRET")
	accessToken, aTokOk := config.Get("TWITTER_ACCESS_TOKEN")
	accessSecret, aSecOk := config.Get("TWITTER_ACCESS_SECRET")
	if !cKeyOk || !cSecOk || !aTokOk || !aSecOk {
		return nil, ErrMissingTokens
	}
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	a.api = anaconda.NewTwitterApi(accessToken, accessSecret)
	return a, nil
}

func (a *API) handler(msg commands.Message, ret commands.MessageFunc) string {
	if len(msg.Params) < 2 {
		return "Usage: .twitter [user]"
	}
	v := url.Values{}
	v.Set("screen_name", msg.Params[1])
	msgs, err := a.api.GetUserTimeline(v)
	if err != nil || len(msgs) == 0 {
		return fmt.Sprintf("I didn't get anything for the %q twitter account", msg.Params[1])
	}
	var time string
	when, err := msgs[0].CreatedAtTime()
	if err != nil {
		time = "Unknown"
	} else {
		time = when.Local().Format("Mon Jan 2, 2006 @ 3:04pm")
	}
	return fmt.Sprintf("%s: %s [%s rt, %s fav] on %s: https://twitter.com/i/web/status/%s", msgs[0].User.ScreenName, util.StripNewLines(msgs[0].Text),
		humanize.Comma(int64(msgs[0].RetweetCount)), humanize.Comma(int64(msgs[0].FavoriteCount)), time, msgs[0].IdStr)
}

// TwitterHandler produces the command handler for Twitter lookups
func (a *API) TwitterHandler() commands.RegisterFunc {
	return func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
		r := regexp.MustCompile("(?i)^\\.tw(itter)?$")
		return twitterTitle, r, a.handler, commands.Options{}
	}
}
