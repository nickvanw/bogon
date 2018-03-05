package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/cmckee-dev/go-alpha-vantage/timeseries"
	humanize "github.com/dustin/go-humanize"
	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
	"golang.org/x/text/message"
)

const intradayTimeFormat = "2006-01-02 15:04:05"

var stockCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.stock$")
	return stockTitle, out, stockLookup, defaultOptions
}

func stockLookup(msg commands.Message, ret commands.MessageFunc) string {
	AV_KEY, ok := config.Get("ALPHAVANTAGE_API")
	if !ok {
		return "I have not been properly configured for this feature"
	}
	client := timeseries.NewClient(AV_KEY)

	resp, err := client.Intraday(msg.Params[1])
	defer resp.Body.Close()
	if err != nil {
		return "Unable to fetch that stock"
	}

	var quotes stockResponseIntra
	if err := json.NewDecoder(resp.Body).Decode(&quotes); err != nil {
		return "The response given was not what I expected"
	}

	loc, err := time.LoadLocation(quotes.MetaData.FiveTimeZone)
	if err != nil {
		return "Error parsing the timezone on that response"
	}

	var data []time.Time
	for k := range quotes.TimeSeries {
		stockTime, _ := time.ParseInLocation(intradayTimeFormat, k, loc)
		data = append(data, stockTime)
	}

	sort.Slice(data, func(i int, j int) bool {
		return data[j].Before(data[i])
	})

	firstQuote := data[0].Format(intradayTimeFormat)
	firstQuotePretty := humanize.Time(data[0])
	firstData := quotes.TimeSeries[firstQuote]

	p := message.NewPrinter(message.MatchLanguage("en"))
	currentInfo := p.Sprintf("Open: %.2f | High: %.2f | Low: %.2f | Close: %.2f",
		firstData.Open, firstData.High, firstData.Low, firstData.Close)

	return fmt.Sprintf("[%s]: %s [quote from %s]",
		quotes.MetaData.TwoSymbol, currentInfo, firstQuotePretty)
}

type stockResponseDaily struct {
	stockMetaData
	TimeSeries map[string]stockQuote `json:"Time Series (Daily)"`
}

type stockResponseIntra struct {
	stockMetaData
	TimeSeries map[string]stockQuote `json:"Time Series (15min)"`
}

type stockMetaData struct {
	MetaData struct {
		OneInformation     string `json:"1. Information"`
		TwoSymbol          string `json:"2. Symbol"`
		ThreeLastRefreshed string `json:"3. Last Refreshed"`
		FourOutputSize     string `json:"4. Output Size"`
		FiveTimeZone       string `json:"5. Time Zone"`
	} `json:"Meta Data"`
}

type stockQuote struct {
	Open   float64 `json:"1. open,string"`
	High   float64 `json:"2. high,string"`
	Low    float64 `json:"3. low,string"`
	Close  float64 `json:"4. close,string"`
	Volume int     `json:"5. volume,string"`
}
