package cmd

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/ChimeraCoder/anaconda"
	"github.com/dustin/go-humanize"
)

var (
	twitterApi *anaconda.TwitterApi
	once       sync.Once
)

func init() {
	AddPlugin("Twitter", "(?i)^\\.tw(itter)?$", MessageHandler(Twitter), false, false)
}

func Twitter(msg *Message) {
	once.Do(initTwitter)
	if twitterApi == nil {
		return
	}
	if len(msg.Params) < 2 {
		msg.Return("Usage: .twitter [user]")
		return
	}
	v := url.Values{}
	v.Set("screen_name", msg.Params[1])
	msgs, err := twitterApi.GetUserTimeline(v)
	if err != nil || len(msgs) == 0 {
		msg.Return("I didn't get anything for " + msg.Params[1] + "'s twitter")
		return
	}
	var time string
	when, err := msgs[0].CreatedAtTime()
	if err != nil {
		time = "Unknown"
	} else {
		time = when.Format("Mon Jan 2, 2006 @ 3:04pm")
	}
	replyMsg := fmt.Sprintf("%s: %s [%s rt, %s fav] on %s", msgs[0].User.ScreenName, stripNewLines(msgs[0].Text),
		humanize.Comma(int64(msgs[0].RetweetCount)), humanize.Comma(int64(msgs[0].FavoriteCount)), time)
	msg.Return(replyMsg)
}

func initTwitter() {
	consumerKey, c_key := GetConfig("TwitterConsumerKey")
	consumerSecret, c_sec := GetConfig("TwitterConsumerSecret")
	accessToken, a_tok := GetConfig("TwitterAccessToken")
	accessSecret, a_sec := GetConfig("TwitterAccessSecret")
	if !c_key || !c_sec || !a_tok || !a_sec {
		return
	}
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	twitterApi = anaconda.NewTwitterApi(accessToken, accessSecret)
}
