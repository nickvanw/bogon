package plugins

import (
	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/plugins/bing"
	"github.com/nickvanw/bogon/commands/plugins/spotify"
	"github.com/nickvanw/bogon/commands/plugins/youtube"
)

// Names of commands exported
const (
	bsTitle        = "bs"
	btcTitle       = "btc"
	currencyTitle  = "currency"
	defineTitle    = "define"
	dnsTitle       = "dns"
	forecastTitle  = "forecast"
	ipTitle        = "ip"
	ltcTitle       = "ltc"
	mehTitle       = "meh"
	redditTitle    = "reddit"
	stockTitle     = "stock"
	subredditTitle = "subreddit"
	tumblrTitle    = "tumblr"
	titpTitle      = "titp"
	urbanTitle     = "urban"
	walkscoreTitle = "walkscore"
	wikiTitle      = "wikipedia"
	weatherTitle   = "weather"
	wolframTitle   = "wolfram"
	waliquorTitle  = "waliquor"
	gifmeTitle     = "gifme"
)

// exports is the list of commands for use in this package
var exports = []commands.RegisterFunc{bsCommand, bitcoinCommand, dnsCommand, ipLookup,
	waLiquor, ltcCommand, currencyCommand, mehCommand, defineCommand, forecastCommand,
	redditCommand, stockCommand, subredditCommand, titpCommand, tumblrCommand, urbanCommand,
	walkscoreCommand, weatherCommand, wikiCommand, wolframCommand, youtube.YoutubeCommand,
	bing.ImageSearch, bing.BingSearch, spotify.Spotify, gifmeCommand}

// Exports is used to return the current registered plugin methods
func Exports() []commands.RegisterFunc {
	return exports
}

// defaultOptions
var defaultOptions = commands.Options{}
