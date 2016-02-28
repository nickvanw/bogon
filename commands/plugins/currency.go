package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/bogon/commands/util"
)

type exchangeData struct {
	sync.RWMutex
	r exchangeRates
}

var data = new(exchangeData)

var currencyCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.(convert)$")
	return currencyTitle, out, convertCurrency, defaultOptions
}

func convertCurrency(msg commands.Message, ret commands.MessageFunc) string {
	if len(msg.Params) < 5 {
		return "Usage: .convert [AMNT] [BASE CURRENCY] (in|to) [CONVERTED CURRENCY], ie: .convert 20 USD in GBP"
	}
	baseRate, err := strconv.ParseFloat(msg.Params[1], 64)
	if err != nil {
		return "I couldn't parse the base rate, type the command with no arguments for help"
	}
	data.Lock()
	defer data.Unlock()
	if time.Duration(time.Now().Sub(time.Unix(data.r.Timestamp, 0))) > 3*time.Hour {
		apikey, avail := config.Get("OPENEXCHANGE_API")
		if avail != true {
			return ""
		}
		if err := updateRates(apikey); err != nil {
			return "I was unable to update my exchange rates, apologies"
		}
	}
	fromCur, fromCurOk := data.r.Rates[strings.ToUpper(msg.Params[2])]
	toCur, toCurOk := data.r.Rates[strings.ToUpper(msg.Params[4])]
	if !toCurOk || !fromCurOk {
		return "I don't have the currencies you wanted listed."
	}
	convAmnt := baseRate * (toCur / fromCur)
	return fmt.Sprintf("%v %s = %.2f %s (%.4f %s per %s)", baseRate, msg.Params[2],
		convAmnt, msg.Params[4], (toCur / fromCur), msg.Params[4], msg.Params[2])
}

func updateRates(apikey string) error {
	url := fmt.Sprintf("http://openexchangerates.org/api/latest.json?app_id=%s", apikey)
	ret, err := util.Fetch(url)
	if err != nil {
		return err
	}
	err = json.Unmarshal(ret, &data.r)
	if err != nil {
		return err
	}
	return nil
}

type exchangeRates struct {
	Base       string             `json:"base"`
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Rates      map[string]float64 `json:"rates"`
	Timestamp  int64              `json:"timestamp"`
}
