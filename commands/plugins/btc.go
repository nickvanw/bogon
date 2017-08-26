package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/util"
)

var bitcoinCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.(btc|bitcoin)$")
	return btcTitle, out, bitCoin, defaultOptions
}

func bitCoin(_ commands.Message, ret commands.MessageFunc) string {
	out := "BTC -> USD:"
	ch := make(chan string)
	go getBitstamp(ch)
	go getCoinbase(ch)
	go getCoinMarket(ch)
	for i := 0; i < 3; i++ {
		val := <-ch
		out += " " + val
	}
	return out
}

func getCoinbase(ch chan string) {
	btce, err := util.Fetch("https://api.coindesk.com/v1/bpi/currentprice.json")
	if err != nil {
		ch <- fmt.Sprintf("[%s]: Coinbase Error!", util.Bold("Coinbase"))
	}

	var response coinbaseResponse
	if err := json.Unmarshal(btce, &response); err != nil {
		ch <- fmt.Sprintf("[%s]: Coinbase Error!", util.Bold("Coinbase"))
		return
	}
	out := fmt.Sprintf("[%s]: Current Rate: $%s", util.Bold("Coinbase"), response.Bpi.USD.Rate)
	ch <- out
}

func getBitstamp(ch chan string) {
	bitstamp, err := util.Fetch("https://www.bitstamp.net/api/ticker/")
	if err != nil {
		ch <- fmt.Sprintf("[%s]: BitStamp Error!", util.Bold("BITSTAMP"))
		return
	}
	var btresponse bitstampResponse
	if err := json.Unmarshal(bitstamp, &btresponse); err != nil {
		ch <- fmt.Sprintf("[%s]: BitStamp Error!", util.Bold("BITSTAMP"))
		return
	}
	out := fmt.Sprintf("[%s]: Last: $%s, High: $%s, Low: $%s", util.Bold("BITSTAMP"), btresponse.Last, btresponse.High, btresponse.Low)
	ch <- out
}

func getCoinMarket(ch chan string) {
	data, err := util.Fetch("https://api.coinmarketcap.com/v1/ticker/bitcoin/")
	if err != nil {
		ch <- fmt.Sprintf("[%s]: CMK Error!", util.Bold("CMK"))
		return
	}
	var response coinMarketCapResponse
	if err := json.Unmarshal(data, &response); err != nil {
		ch <- fmt.Sprintf("[%s]: CMK Error!", util.Bold("CMK"))
		return
	}
	out := fmt.Sprintf("[%s] Current Rate: $%s", util.Bold("CMK"), response[0].PriceUsd)
	ch <- out
}

type bitstampResponse struct {
	Ask       string `json:"ask"`
	Bid       string `json:"bid"`
	High      string `json:"high"`
	Last      string `json:"last"`
	Low       string `json:"low"`
	Timestamp string `json:"timestamp"`
	Volume    string `json:"volume"`
}

type coinbaseResponse struct {
	Bpi struct {
		USD struct {
			Code   string `json:"code"`
			Symbol string `json:"symbol"`
			Rate   string `json:"rate"`
		} `json:"USD"`
	} `json:"bpi"`
}

type coinMarketCapResponse []struct {
	PriceUsd string `json:"price_usd"`
}
