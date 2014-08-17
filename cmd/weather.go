package cmd

import (
	"encoding/json"
	"fmt"
)

func init() {
	AddPlugin("Weather", "(?i)^\\.weather$", MessageHandler(Weather), false, false)
}

type Conditions struct {
	Current_observation Current
}

type Current struct {
	Observation_location Location
	Station_id           string
	Weather              string
	Temperature_string   string
	Relative_humidity    string
	Wind_string          string
	Precip_today_string  string
}

type Location struct {
	Full string
}

func Weather(msg *Message) {
	geoAddr, err := GetCoordinates(msg.Params[1:])
	if err != nil {
		msg.Return("I couldn't track down that location!")
		return
	}
	api_key, avail := GetConfig("Wunderground")
	if avail != true {
		fmt.Println("wunderground API Key not available!")
		return
	}
	url := fmt.Sprintf("http://api.wunderground.com/api/%s/conditions/q/%v,%v.json", api_key, geoAddr.Lat, geoAddr.Long)
	data, err := getSite(url)
	if err != nil {
		msg.Return("Error!")
		return
	}
	var conditions Conditions
	json.Unmarshal(data, &conditions)
	response := conditions.Current_observation
	if response.Weather == "" {
		msg.Return("I couldn't find that location -- try again!")
		return
	}
	location := fmt.Sprintf("%s (%s)", response.Observation_location.Full, response.Station_id)
	out := fmt.Sprintf("%s is: %s - %s with %s humidity | Wind: %s | %s precip. today", location, response.Weather, response.Temperature_string, response.Relative_humidity, response.Wind_string, response.Precip_today_string)
	msg.Return(out)
}
