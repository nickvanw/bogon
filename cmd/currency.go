package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ExchangeData ExchangeRates

func init() {
	AddPlugin("Convert", "(?i)^\\.(convert)$", MessageHandler(Convert), false, false)
}

func Convert(msg *Message) {
	if len(msg.Params) < 5 {
		msg.Return("Usage: .convert [AMNT] [BASE CURRENCY] (in|to) [CONVERTED CURRENCY], ie: .convert 20 USD in GBP")
		return
	}
	base_rate, err := strconv.ParseFloat(msg.Params[1], 64)
	if err != nil {
		msg.Return("I couldn't parse the base rate, type the command with no arguments for help")
	}
	if time.Duration(time.Now().Sub(time.Unix(ExchangeData.Timestamp, 0))) > time.Hour*3 {
		apikey, avail := GetConfig("Convert")
		if avail != true {
			fmt.Println("I don't have an API Key for openexchangerates.org")
			return
		}
		success := UpdateExchangeRates(apikey)
		if !success {
			msg.Return("I was unable to update my exchange rates, apologies")
			return
		}
	}
	from_cur, from_cur_ok := ExchangeData.Rates[strings.ToUpper(msg.Params[2])]
	to_cur, to_cur_ok := ExchangeData.Rates[strings.ToUpper(msg.Params[4])]
	if !to_cur_ok || !from_cur_ok {
		msg.Return("I don't have the currencies you wanted listed. For more information visit: http://sys.nvw.io/cur.html")
		return
	}
	converted_amnt := base_rate * (to_cur / from_cur)
	msg.Return(fmt.Sprintf("%v %s = %.2f %s (%.4f %s per %s)", base_rate, msg.Params[2], converted_amnt, msg.Params[4], (to_cur / from_cur), msg.Params[4], msg.Params[2]))
}

func UpdateExchangeRates(apikey string) bool {
	url := fmt.Sprintf("http://openexchangerates.org/api/latest.json?app_id=%s", apikey)
	data, err := getSite(url)
	if err != nil {
		return false
	}
	err = json.Unmarshal(data, &ExchangeData)
	if err != nil {
		return false
	}
	return true
}

type ExchangeRates struct {
	Base       string             `json:"base"`
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Rates      map[string]float64 `json:"rates"`
	Timestamp  int64              `json:"timestamp"`
}
