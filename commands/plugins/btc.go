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
	go getBTCE(ch)
	for i := 0; i < 2; i++ {
		val := <-ch
		out += " " + val
	}
	return out
}

func getBTCE(ch chan string) {
	btce, err := util.Fetch("https://btc-e.com/api/2/btc_usd/ticker")
	if err != nil {
		ch <- fmt.Sprintf("[%s]: BTC-E Error!", util.Bold("BTC-E"))
	}

	var response btceResponse
	if err := json.Unmarshal(btce, &response); err != nil {
		ch <- fmt.Sprintf("[%s]: BTC-E Error!", util.Bold("BTC-E"))
		return
	}
	out := fmt.Sprintf("[%s]: Last: $%v, High: $%v, Low: $%v, Avg: $%v", util.Bold("BTC-E"), response.Ticker.Last, response.Ticker.High, response.Ticker.Low, response.Ticker.Avg)
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

type bitstampResponse struct {
	Ask       string `json:"ask"`
	Bid       string `json:"bid"`
	High      string `json:"high"`
	Last      string `json:"last"`
	Low       string `json:"low"`
	Timestamp string `json:"timestamp"`
	Volume    string `json:"volume"`
}

type btceResponse struct {
	Ticker struct {
		Avg        float64 `json:"avg"`
		Buy        float64 `json:"buy"`
		High       float64 `json:"high"`
		Last       float64 `json:"last"`
		Low        float64 `json:"low"`
		Sell       float64 `json:"sell"`
		ServerTime float64 `json:"server_time"`
		Updated    float64 `json:"updated"`
		Vol        float64 `json:"vol"`
		VolCur     float64 `json:"vol_cur"`
	} `json:"ticker"`
}
