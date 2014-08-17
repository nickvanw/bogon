package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func init() {
	AddPlugin("WalkScore", "(?i)^\\.walk(score)?$", MessageHandler(WalkScore), false, false)
}

func WalkScore(msg *Message) {
	walkscore, avail := GetConfig("WalkScore")
	if avail != true {
		fmt.Println("WalkScore API Key Not available")
		return
	}
	geoAddr, err := GetCoordinates(msg.Params[1:])
	if err != nil {
		msg.Return("Not found!")
		return
	}
	lat := geoAddr.Lat
	long := geoAddr.Long
	addr := url.QueryEscape(geoAddr.FormattedAddress)
	walkscoreURL := fmt.Sprintf("http://api.walkscore.com/score?format=json&address=%v&lat=%v&lon=%v&wsapikey=%v", addr, lat, long, walkscore)
	wdata, err := getSite(walkscoreURL)
	if err != nil {
		msg.Return("Error fetching walkscore data")
		return
	}
	var ws WalkScoreData
	json.Unmarshal(wdata, &ws)
	fmt.Println(string(wdata))
	if ws.Status != 1 {
		msg.Return("I couldn't find that in Walkscore's database")
		return
	}
	out := fmt.Sprintf("%s is a %s with a walkscore of %v", geoAddr.FormattedAddress, ws.Description, ws.Walkscore)
	msg.Return(out)
}

func GetCoordinates(addr []string) (GoogleReturn, error) {
	address := urlencode(strings.Join(addr, " "))
	geoURL := fmt.Sprintf("http://maps.googleapis.com/maps/api/geocode/json?address=%s&sensor=false", address)
	data, err := getSite(geoURL)
	if err != nil {
		return GoogleReturn{}, errors.New("Invalid Address")
	}
	var geo Geolocation
	json.Unmarshal(data, &geo)
	if geo.Status != "OK" || len(geo.Results) < 1 {
		return GoogleReturn{}, errors.New("Invalid Address")
	}
	ret := GoogleReturn{
		Lat:              geo.Results[0].Geometry.Location.Lat,
		Long:             geo.Results[0].Geometry.Location.Lng,
		FormattedAddress: geo.Results[0].FormattedAddress,
	}
	return ret, nil
}

type GoogleReturn struct {
	Lat              float64
	Long             float64
	FormattedAddress string
}

type WalkScoreData struct {
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

type Geolocation struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		Types []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}
