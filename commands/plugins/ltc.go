package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/util"
)

var ltcCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.(ltc|litecoin)$")
	return ltcTitle, out, getLtc, defaultOptions
}

func getLtc(_ commands.Message, ret commands.MessageFunc) string {
	btce, err := util.Fetch("https://btc-e.com/api/2/ltc_usd/ticker")
	if err != nil {
		return fmt.Sprintf("[%s]: BTC-E Error!", util.Bold("BTC-E"))
	}

	var response btceResponse
	if err := json.Unmarshal(btce, &response); err != nil {
		return fmt.Sprintf("[%s]: BTC-E Error!", util.Bold("BTC-E"))
	}
	return fmt.Sprintf("LTC->USD: Last: $%v, High: $%v, Low: $%v, Avg: $%v",
		response.Ticker.Last, response.Ticker.High, response.Ticker.Low, response.Ticker.Avg)
}
