package plugins

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/bogon/commands/util"
)

var walkscoreCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.walk(score)?$")
	return walkscoreTitle, out, walkscoreLookup, defaultOptions
}

func walkscoreLookup(msg commands.Message, ret commands.MessageFunc) string {
	walkscore, avail := config.Get("WALKSCORE_API")
	if avail != true {
		return ""
	}
	geoAddr, err := util.GetCoordinates(msg.Params[1:])
	if err != nil {
		return "Not found!"
	}
	lat := geoAddr.Lat
	long := geoAddr.Long
	addr := url.QueryEscape(geoAddr.FormattedAddress)
	walkscoreURL := fmt.Sprintf("http://api.walkscore.com/score?format=json&address=%v&lat=%v&lon=%v&wsapikey=%v", addr, lat, long, walkscore)
	wdata, err := util.Fetch(walkscoreURL)
	if err != nil {
		return "Error fetching walkscore data"
	}
	var ws walkScoreData
	if err := json.Unmarshal(wdata, &ws); err != nil {
		return "Walkscore gave me a bad response!"
	}
	if ws.Status != 1 {
		return "I couldn't find that in Walkscore's database"
	}
	return fmt.Sprintf("%s is a %s with a walkscore of %v", geoAddr.FormattedAddress, ws.Description, ws.Walkscore)
}

type walkScoreData struct {
	Description  string  `json:"description"`
	LogoURL      string  `json:"logo_url"`
	MoreInfoIcon string  `json:"more_info_icon"`
	MoreInfoLink string  `json:"more_info_link"`
	SnappedLat   float64 `json:"snapped_lat"`
	SnappedLon   float64 `json:"snapped_lon"`
	Status       float64 `json:"status"`
	Updated      string  `json:"updated"`
	Walkscore    float64 `json:"walkscore"`
	WsLink       string  `json:"ws_link"`
}
