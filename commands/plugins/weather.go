package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/bogon/commands/util"
)

var weatherCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.weather?$")
	return weatherTitle, out, weatherLookup, defaultOptions
}

func weatherLookup(msg commands.Message, ret commands.MessageFunc) string {
	geoAddr, err := util.GetCoordinates(msg.Params[1:])
	if err != nil {
		return err.Error()
	}
	apiKey, avail := config.Get("WUNDERGROUND_API")
	if avail != true {
		return ""
	}
	url := fmt.Sprintf("http://api.wunderground.com/api/%s/conditions/q/%v,%v.json", apiKey, geoAddr.Lat, geoAddr.Long)
	data, err := util.Fetch(url)
	if err != nil {
		return "Unable to lookup weather"
	}
	var conditions weatherConditions
	if err := json.Unmarshal(data, &conditions); err != nil {
		return "Weather Underground gave me a bad response"
	}
	response := conditions.Current
	if response.Weather == "" {
		return "I couldn't find that location -- try again!"
	}
	location := fmt.Sprintf("%s (%s)", response.Location.Full, response.StationID)
	return fmt.Sprintf("%s is: %s - %s with %s humidity | Wind: %s | %s precip. today",
		location, response.Weather, response.TempString, response.RelHumidity,
		response.WindString, response.PrecipString)
}

type weatherConditions struct {
	Current current `json:"current_observation"`
}

type current struct {
	Location     location `json:"observation_location"`
	StationID    string   `json:"station_id"`
	Weather      string   `json:"weather"`
	TempString   string   `json:"temperature_string"`
	RelHumidity  string   `json:"relative_humidity"`
	WindString   string   `json:"wind_string"`
	PrecipString string   `json:"precip_today_string"`
}

type location struct {
	Full string `json:"full"`
}
