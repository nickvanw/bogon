package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/util"
)

var ltcCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.(ltc|litecoin)$")
	return ltcTitle, out, getLtc, defaultOptions
}

func getLtc(_ commands.Message, ret commands.MessageFunc) string {
	coinMarket, err := util.Fetch("https://api.coinmarketcap.com/v1/ticker/litecoin/")
	if err != nil {
		return fmt.Sprintf("[%s]: CMK Error!", util.Bold("CMK"))
	}

	var response coinMarketCapResponse
	if err := json.Unmarshal(coinMarket, &response); err != nil {
		return fmt.Sprintf("[%s]: CMK Error!", util.Bold("CMK"))
	}
	return fmt.Sprintf("1 LTC in USD: Current: $%s (%s%% 1w change) ", response[0].PriceUsd, response[0].WeekChange)
}

var ethCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.(eth|ethereum)$")
	return ethTitle, out, getEth, defaultOptions
}

func getEth(_ commands.Message, ret commands.MessageFunc) string {
	coinMarket, err := util.Fetch("https://api.coinmarketcap.com/v1/ticker/ethereum/")
	if err != nil {
		return fmt.Sprintf("[%s]: CMK Error!", util.Bold("CMK"))
	}

	var response coinMarketCapResponse
	if err := json.Unmarshal(coinMarket, &response); err != nil {
		return fmt.Sprintf("[%s]: CMK Error!", util.Bold("CMK"))
	}
	return fmt.Sprintf("1 ETH in USD: Current: $%s (%s%% 1w change)", response[0].PriceUsd, response[0].WeekChange)
}

var dogeCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.(doge|dogecoin)$")
	return ethTitle, out, getDoge, defaultOptions
}

func getDoge(_ commands.Message, ret commands.MessageFunc) string {
	coinMarket, err := util.Fetch("https://api.coinmarketcap.com/v1/ticker/dogecoin/")
	if err != nil {
		return fmt.Sprintf("[%s]: CMK Error!", util.Bold("CMK"))
	}

	var response coinMarketCapResponse
	if err = json.Unmarshal(coinMarket, &response); err != nil {
		return fmt.Sprintf("[%s]: CMK Error!", util.Bold("CMK"))
	}
	price, err := strconv.ParseFloat(response[0].PriceUsd, 64)
	if err != nil {
		return fmt.Sprintf("[%s]: CMK Error!", util.Bold("CMK"))
	}
	usdPrice := 1 / price

	return fmt.Sprintf("1 Dogecoin in USD: $%s | 1 USD = %.2f Doge (%s%% 1w change)", response[0].PriceUsd, usdPrice, response[0].WeekChange)
}
