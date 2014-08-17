package cmd

import (
	"encoding/json"
	"fmt"
)

func init() {
	AddPlugin("Litecoin", "(?i)^\\.(ltc|litecoin)$", MessageHandler(LiteCoin), false, false)
}

func LiteCoin(msg *Message) {
	data := getLtcE()
	msg.Return(data)
}

func getLtcE() string {
	btce, err := getSite("https://btc-e.com/api/2/ltc_usd/ticker")
	if err != nil {
		return fmt.Sprintf("[%s]: BTC-E Error!", bold("BTC-E"))
	}

	var response BTCE
	json.Unmarshal(btce, &response)
	out := fmt.Sprintf("LTC->USD: Last: $%v, High: $%v, Low: $%v, Avg: $%v", response.Ticker.Last, response.Ticker.High, response.Ticker.Low, response.Ticker.Avg)
	return out
}
