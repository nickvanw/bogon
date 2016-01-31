package plugins

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	forecast "github.com/mlbright/forecast/v2"
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

	f, err := forecast.Get(apikey, lat, long, "now", forecast.US)
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
