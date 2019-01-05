package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/util"
	"golang.org/x/text/message"
)

const apiURL = "https://api.iextrading.com/1.0/stock/%s/book"

var stockCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.stock$")
	return stockTitle, out, stockLookup, defaultOptions
}

func stockLookup(msg commands.Message, ret commands.MessageFunc) string {
	dataURL := fmt.Sprintf(apiURL, msg.Params[1])

	data, err := util.Fetch(dataURL)
	if err != nil {
		return err.Error()
	}

	var response stockQuote
	if err := json.Unmarshal(data, &response); err != nil {
		return "I got invalid data for that quote"
	}

	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprintf("%s (%s): Latest Price: %.2f (from '%s' at %s) | Open: %.2f | Close: %.2f | Change: %.2f (%.2f%%) | 52 Week High: %.2f | 52 Week Low: %.2f | YTD Change: %.2f%%",
		response.Quote.Symbol, response.Quote.CompanyName, response.Quote.LatestPrice,
		response.Quote.LatestSource, response.Quote.LatestTime, response.Quote.Open,
		response.Quote.Close, response.Quote.Change, response.Quote.ChangePercent*100,
		response.Quote.Week52High, response.Quote.Week52Low, response.Quote.YtdChange*100)

}

type stockQuote struct {
	Quote struct {
		Symbol           string  `json:"symbol"`
		CompanyName      string  `json:"companyName"`
		PrimaryExchange  string  `json:"primaryExchange"`
		Sector           string  `json:"sector"`
		CalculationPrice string  `json:"calculationPrice"`
		Open             float64 `json:"open"`
		OpenTime         int64   `json:"openTime"`
		Close            float64 `json:"close"`
		CloseTime        int64   `json:"closeTime"`
		High             float64 `json:"high"`
		Low              float64 `json:"low"`
		LatestPrice      float64 `json:"latestPrice"`
		LatestSource     string  `json:"latestSource"`
		LatestTime       string  `json:"latestTime"`
		LatestUpdate     int64   `json:"latestUpdate"`
		LatestVolume     int     `json:"latestVolume"`
		IexRealtimePrice float64 `json:"iexRealtimePrice"`
		IexRealtimeSize  int     `json:"iexRealtimeSize"`
		IexLastUpdated   int64   `json:"iexLastUpdated"`
		DelayedPrice     float64 `json:"delayedPrice"`
		DelayedPriceTime int64   `json:"delayedPriceTime"`
		PreviousClose    float64 `json:"previousClose"`
		Change           float64 `json:"change"`
		ChangePercent    float64 `json:"changePercent"`
		IexMarketPercent float64 `json:"iexMarketPercent"`
		IexVolume        int     `json:"iexVolume"`
		AvgTotalVolume   int     `json:"avgTotalVolume"`
		IexBidPrice      float64 `json:"iexBidPrice"`
		IexBidSize       int     `json:"iexBidSize"`
		IexAskPrice      float64 `json:"iexAskPrice"`
		IexAskSize       int     `json:"iexAskSize"`
		MarketCap        int64   `json:"marketCap"`
		PeRatio          float64 `json:"peRatio"`
		Week52High       float64 `json:"week52High"`
		Week52Low        float64 `json:"week52Low"`
		YtdChange        float64 `json:"ytdChange"`
	} `json:"quote"`
	SystemEvent struct {
		SystemEvent string `json:"systemEvent"`
		Timestamp   int64  `json:"timestamp"`
	} `json:"systemEvent"`
}
