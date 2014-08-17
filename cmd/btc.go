package cmd

import (
	"encoding/json"
	"fmt"
)

func init() {
	AddPlugin("Bitcoin", "(?i)^\\.(btc|bitcoin)$", MessageHandler(BitCoin), false, false)
}

func BitCoin(msg *Message) {
	out := "BTC -> USD:"
	ch := make(chan string)
	go getBitStamp(ch)
	go getBtcE(ch)
	for i := 0; i < 2; i++ {
		val := <-ch
		out = fmt.Sprintf("%s %s", out, val)
	}
	msg.Return(out)
}

func getBtcE(ch chan string) {
	btce, err := getSite("https://btc-e.com/api/2/btc_usd/ticker")
	if err != nil {
		ch <- fmt.Sprintf("[%s]: BTC-E Error!", bold("BTC-E"))
	}

	var response BTCE
	json.Unmarshal(btce, &response)
	out := fmt.Sprintf("[%s]: Last: $%v, High: $%v, Low: $%v, Avg: $%v", bold("BTC-E"), response.Ticker.Last, response.Ticker.High, response.Ticker.Low, response.Ticker.Avg)
	ch <- out
}

func getBitStamp(ch chan string) {
	bitstamp, err := getSite("https://www.bitstamp.net/api/ticker/")
	if err != nil {
		ch <- fmt.Sprintf("[%s]: BitStamp Error!", bold("BITSTAMP"))
	}
	var btresponse BSResponse
	json.Unmarshal(bitstamp, &btresponse)
	out := fmt.Sprintf("[%s]: Last: $%s, High: $%s, Low: $%s", bold("BITSTAMP"), btresponse.Last, btresponse.High, btresponse.Low)
	ch <- out
}

type MTGoxReturn struct {
	Result string
	Return map[string]map[string]string
}

type BSResponse struct {
	Ask       string `json:"ask"`
	Bid       string `json:"bid"`
	High      string `json:"high"`
	Last      string `json:"last"`
	Low       string `json:"low"`
	Timestamp string `json:"timestamp"`
	Volume    string `json:"volume"`
}

type BTCE struct {
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
