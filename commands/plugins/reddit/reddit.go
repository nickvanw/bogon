package reddit

import (
	"regexp"

	"github.com/nickvanw/bogon/commands"
)

var (
	userUrl = "https://oauth.reddit.com/user/%s/about.json"
	subUrl  = "https://oauth.reddit.com/r/%s.json"
)

const (
	redditUserTitle = "reddit_user"
	redditSubTitle  = "reddit_sub"
)

var RedditUser = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.r(eddit)?u(ser)?$")
	return redditUserTitle, out, redditUserHandler, commands.Options{}
}

var RedditSub = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.(sub)?r(eddit)?$")
	return redditSubTitle, out, redditSubHandler, commands.Options{}
}

func redditUserHandler(msg commands.Message, ret commands.MessageFunc) string {
	if len(msg.Params) < 2 {
		return "Usage .ru [username]"
	}
	user := msg.Params[1]
	return userLookup(user)
}

func redditSubHandler(msg commands.Message, ret commands.MessageFunc) string {
	if len(msg.Params) < 2 {
		return "Usage .r [subreddit]"
	}
	sub := msg.Params[1]
	return subLookup(sub)
}
