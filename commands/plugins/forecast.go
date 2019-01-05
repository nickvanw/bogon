package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/bogon/commands/util"
)

var forecastCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.forecast$")
	return forecastTitle, out, fetchForecast, defaultOptions
}

func fetchForecast(msg commands.Message, ret commands.MessageFunc) string {
	geoAddr, err := util.GetCoordinates(msg.Params[1:])
	if err != nil {
		return "I couldn't track down that location!"

	}
	apikey, avail := config.Get("FORECASTIO_API")
	if avail != true {
		return ""
	}

	lat := strconv.FormatFloat(geoAddr.Lat, 'f', -1, 64)
	long := strconv.FormatFloat(geoAddr.Long, 'f', -1, 64)

	f, err := getForecast(apikey, lat, long, "now", US, English)
	if err != nil {
		return "Unable to fetch forecast!"
	}
	out := fmt.Sprintf("Forecast: %s", f.Daily.Summary)
	for i := 0; i < int(math.Min(float64(len(f.Daily.Data)), 3.0)); i++ {
		data := f.Daily.Data[i]
		out = fmt.Sprintf("%s | %s: %s [high: %.0fF, low: %.0fF]", out,
			time.Unix(int64(data.Time), 0).Format("Mon Jan 2"), data.Summary, data.TemperatureMax, data.TemperatureMin)
		if data.PrecipType != "" {
			out = fmt.Sprintf("%s; [%.0f%% of %s]", out, data.PrecipProbability*100.0, data.PrecipType)
		}
	}
	return out
}

// URL example:  "https://api.darksky.net/forecast/APIKEY/LATITUDE,LONGITUDE,TIME?units=ca&lang=en"
const (
	BASEURL = "https://api.darksky.net/forecast"
)

type Flags struct {
	DarkSkyUnavailable string   `json:"darksky-unavailable,omitempty"`
	DarkSkyStations    []string `json:"darksky-stations,omitempty"`
	DataPointStations  []string `json:"datapoint-stations,omitempty"`
	ISDStations        []string `json:"isds-stations,omitempty"`
	LAMPStations       []string `json:"lamp-stations,omitempty"`
	MADISStations      []string `json:"madis-stations,omitempty"`
	METARStations      []string `json:"metars-stations,omitempty"`
	METNOLicense       string   `json:"metnol-license,omitempty"`
	Sources            []string `json:"sources,omitempty"`
	Units              string   `json:"units,omitempty"`
}

type DataPoint struct {
	Time                       int64   `json:"time,omitempty"`
	Summary                    string  `json:"summary,omitempty"`
	Icon                       string  `json:"icon,omitempty"`
	SunriseTime                int64   `json:"sunriseTime,omitempty"`
	SunsetTime                 int64   `json:"sunsetTime,omitempty"`
	PrecipIntensity            float64 `json:"precipIntensity,omitempty"`
	PrecipIntensityMax         float64 `json:"precipIntensityMax,omitempty"`
	PrecipIntensityMaxTime     int64   `json:"precipIntensityMaxTime,omitempty"`
	PrecipProbability          float64 `json:"precipProbability,omitempty"`
	PrecipType                 string  `json:"precipType,omitempty"`
	PrecipAccumulation         float64 `json:"precipAccumulation,omitempty"`
	Temperature                float64 `json:"temperature,omitempty"`
	TemperatureMin             float64 `json:"temperatureMin,omitempty"`
	TemperatureMinTime         int64   `json:"temperatureMinTime,omitempty"`
	TemperatureMax             float64 `json:"temperatureMax,omitempty"`
	TemperatureMaxTime         int64   `json:"temperatureMaxTime,omitempty"`
	ApparentTemperature        float64 `json:"apparentTemperature,omitempty"`
	ApparentTemperatureMin     float64 `json:"apparentTemperatureMin,omitempty"`
	ApparentTemperatureMinTime int64   `json:"apparentTemperatureMinTime,omitempty"`
	ApparentTemperatureMax     float64 `json:"apparentTemperatureMax,omitempty"`
	ApparentTemperatureMaxTime int64   `json:"apparentTemperatureMaxTime,omitempty"`
	NearestStormBearing        float64 `json:"nearestStormBearing,omitempty"`
	NearestStormDistance       float64 `json:"nearestStormDistance,omitempty"`
	DewPoint                   float64 `json:"dewPoint,omitempty"`
	WindSpeed                  float64 `json:"windSpeed,omitempty"`
	WindBearing                float64 `json:"windBearing,omitempty"`
	CloudCover                 float64 `json:"cloudCover,omitempty"`
	Humidity                   float64 `json:"humidity,omitempty"`
	Pressure                   float64 `json:"pressure,omitempty"`
	Visibility                 float64 `json:"visibility,omitempty"`
	Ozone                      float64 `json:"ozone,omitempty"`
	MoonPhase                  float64 `json:"moonPhase,omitempty"`
}

type dataBlock struct {
	Summary string      `json:"summary,omitempty"`
	Icon    string      `json:"icon,omitempty"`
	Data    []DataPoint `json:"data,omitempty"`
}

type alert struct {
	Title       string   `json:"title,omitempty"`
	Regions     []string `json:"regions,omitempty"`
	Severity    string   `json:"severity,omitempty"`
	Description string   `json:"description,omitempty"`
	Time        int64    `json:"time,omitempty"`
	Expires     float64  `json:"expires,omitempty"`
	URI         string   `json:"uri,omitempty"`
}

type forecast struct {
	Latitude  float64   `json:"latitude,omitempty"`
	Longitude float64   `json:"longitude,omitempty"`
	Timezone  string    `json:"timezone,omitempty"`
	Offset    float64   `json:"offset,omitempty"`
	Currently DataPoint `json:"currently,omitempty"`
	Minutely  dataBlock `json:"minutely,omitempty"`
	Hourly    dataBlock `json:"hourly,omitempty"`
	Daily     dataBlock `json:"daily,omitempty"`
	Alerts    []alert   `json:"alerts,omitempty"`
	Flags     Flags     `json:"flags,omitempty"`
	APICalls  int       `json:"apicalls,omitempty"`
	Code      int       `json:"code,omitempty"`
}

type units string

const (
	CA   units = "ca"
	SI   units = "si"
	US   units = "us"
	UK   units = "uk"
	AUTO units = "auto"
)

type lang string

const (
	Arabic             lang = "ar"
	Azerbaijani        lang = "az"
	Belarusian         lang = "be"
	Bosnian            lang = "bs"
	Catalan            lang = "ca"
	Czech              lang = "cs"
	German             lang = "de"
	Greek              lang = "el"
	English            lang = "en"
	Spanish            lang = "es"
	Estonian           lang = "et"
	French             lang = "fr"
	Croatian           lang = "hr"
	Hungarian          lang = "hu"
	Indonesian         lang = "id"
	Italian            lang = "it"
	Icelandic          lang = "is"
	Cornish            lang = "kw"
	Indonesia          lang = "nb"
	Dutch              lang = "nl"
	Polish             lang = "pl"
	Portuguese         lang = "pt"
	Russian            lang = "ru"
	Slovak             lang = "sk"
	Slovenian          lang = "sl"
	Serbian            lang = "sr"
	Swedish            lang = "sv"
	Tetum              lang = "te"
	Turkish            lang = "tr"
	Ukrainian          lang = "uk"
	IgpayAtinlay       lang = "x-pig-latin"
	SimplifiedChinese  lang = "zh"
	TraditionalChinese lang = "zh-tw"
)

func getForecast(key string, lat string, long string, time string, units units, lang lang) (*forecast, error) {
	res, err := getForecastResponse(key, lat, long, time, units, lang)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	f, err := fromJson(res.Body)
	if err != nil {
		return nil, err
	}

	calls, _ := strconv.Atoi(res.Header.Get("X-Forecast-API-Calls"))
	f.APICalls = calls

	return f, nil
}

func fromJson(reader io.Reader) (*forecast, error) {
	var f forecast
	if err := json.NewDecoder(reader).Decode(&f); err != nil {
		return nil, err
	}

	return &f, nil
}

func getForecastResponse(key string, lat string, long string, time string, units units, lang lang) (*http.Response, error) {
	coord := lat + "," + long

	var url string
	if time == "now" {
		url = BASEURL + "/" + key + "/" + coord + "?units=" + string(units) + "&lang=" + string(lang)
	} else {
		url = BASEURL + "/" + key + "/" + coord + "," + time + "?units=" + string(units) + "&lang=" + string(lang)
	}

	res, err := http.Get(url)
	if err != nil {
		return res, err
	}

	return res, nil
}
