package spotify

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nickvanw/bogon/commands"
)

var (
	spotifyURL  = regexp.MustCompile(`(?i)(spotify:(?:(?:artist|album|track|user:[^:]+:playlist):[a-zA-Z0-9]+))`)
	spotifyLink = regexp.MustCompile(`(?i)http(?:s)?:\/\/(?:open|play)\.spotify\.com\/(track|album|artist)\/([\w]+)`)
	lookupURL   = "https://api.spotify.com/v1"
)

const (
	spotifyTitle = "spotify"
)

// Spotify registers a handler to look up spotify URLs
var Spotify = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	return spotifyTitle, nil, spotifyHandler, commands.Options{Raw: true}
}

func spotifyHandler(msg commands.Message, ret commands.MessageFunc) string {
	var matches []string
	for _, v := range msg.Params {
		if urlMatch := spotifyURL.FindAllStringSubmatch(v, -1); len(urlMatch) > 0 {
			matches = append(matches, urlMatch[0][0])
		}
		if linkMatch := spotifyLink.FindAllStringSubmatch(v, -1); len(linkMatch) > 0 {
			matches = append(matches, fmt.Sprintf("spotify:%s:%s", linkMatch[0][1], linkMatch[0][2]))
		}
	}
	if len(matches) > 0 {
		return processSpotify(matches)
	}
	return ""
}

func processSpotify(matches []string) string {
	pieces := strings.Split(matches[0], ":")
	var out string
	switch strings.ToLower(pieces[1]) {
	case "artist":
		out = artistLookup(pieces[2])
	case "album":
		out = albumLookup(pieces[2])
	case "track":
		out = trackLookup(pieces[2])
	default:
		return ""
	}
	if len(matches) > 1 {
		out = fmt.Sprintf("%s and %d more..", out, len(matches)-1)
	}
	return out
}
