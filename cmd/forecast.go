package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func init() {
	AddPlugin("Forecast", "(?i)^\\.forecast$", MessageHandler(GetForecast), false, false)
}

func GetForecast(msg *Message) {
	location := url.QueryEscape(strings.Join(msg.Params[1:], " "))
	WUNDERGROUND, avail := GetConfig("Wunderground")
	if avail != true {
		fmt.Println("wunderground API key not available")
		return
	}
	url := fmt.Sprintf("http://api.wunderground.com/api/%s/forecast/q/%s.json", WUNDERGROUND, location)
	data, err := getSite(url)
	if err != nil {
		msg.Return("Error!")
		return
	}
	var forecast ForecastConditions
	json.Unmarshal(data, &forecast)
	out := make([]string, 0, 3)
	for _, v := range forecast.Forecast.Simpleforecast.Forecastday {
		day := fmt.Sprintf("[%v] High: %v°F (%v°C) Low: %v°F (%v°C). %v (%v%% chance precip) with %v mph winds and %v%% humidity", v.Date.Weekday_short, v.High["fahrenheit"], v.High["celsius"], v.Low["fahrenheit"], v.Low["celsius"], v.Conditions, v.Pop, v.Avewind.Mph, v.Avehumidity)
		out = append(out, day)
	}
	msg.Return(strings.Join(out[:3], " | "))
}

type ForecastConditions struct {
	Forecast Forecast
}

type Forecast struct {
	Simpleforecast Simpleforecast
}

type Simpleforecast struct {
	Forecastday []Forecastday
}

type Forecastday struct {
	Date        ForecastDate
	High        map[string]string
	Low         map[string]string
	Conditions  string
	Pop         int
	Avewind     ForecastWind
	Avehumidity int
}

type ForecastDate struct {
	Weekday_short string
}

type ForecastWind struct {
	Mph     int
	Kph     int
	Dir     string
	Degrees int
}
