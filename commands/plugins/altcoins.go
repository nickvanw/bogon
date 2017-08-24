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
	coinMarket, err := util.Fetch("https://api.coinmarketcap.com/v1/ticker/litecoin/")
	if err != nil {
		return fmt.Sprintf("[%s]: CMK Error!", util.Bold("CMK"))
	}

	var response coinMarketCapResponse
	if err := json.Unmarshal(coinMarket, &response); err != nil {
		return fmt.Sprintf("[%s]: CMK Error!", util.Bold("CMK"))
	}
	return fmt.Sprintf("1 LTC in USD: Current: $%s", response[0].PriceUsd)
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
	return fmt.Sprintf("1 ETH in USD: Current: $%s", response[0].PriceUsd)
}
