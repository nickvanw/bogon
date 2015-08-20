package cmd

import (
	"fmt"
	"math"
	"strconv"
	"time"

	forecast "github.com/mlbright/forecast/v2"
)

func init() {
	AddPlugin("Forecast", "(?i)^\\.forecast$", MessageHandler(GetForecast), false, false)
}

func GetForecast(msg *Message) {
	geoAddr, err := GetCoordinates(msg.Params[1:])
	if err != nil {
		msg.Return("I couldn't track down that location!")
		return
	}
	apikey, avail := GetConfig("Forecast")
	if avail != true {
		fmt.Println("forecast.io API key not available")
		return
	}

	lat := strconv.FormatFloat(geoAddr.Lat, 'f', -1, 64)
	long := strconv.FormatFloat(geoAddr.Long, 'f', -1, 64)

	f, err := forecast.Get(apikey, lat, long, "now", forecast.US)
	if err != nil {
		msg.Return("Unable to fetch forecast!")
		return
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
	msg.Return(out)
}
