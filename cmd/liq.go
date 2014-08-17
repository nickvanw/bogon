package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

func init() {
	AddPlugin("WALiquor", "(?i)^\\.wal(iquor)?$", MessageHandler(Liq), false, false)
}

func Liq(msg *Message) {
	if len(msg.Params) < 2 {
		msg.Return("Usage: .waliquor [size] [amount]")
		return
	}
	dollar, err := strconv.ParseFloat(strings.Replace(msg.Params[len(msg.Params)-1], "$", "", -1), 64)
	if err != nil {
		msg.Return("I couldn't parse the currency amount | Usage: .waliquor [size] [amount]")
		return
	}
	toLower := strings.ToLower(msg.Params[1])
	mult := 1.0
	if strings.HasSuffix(toLower, "ml") || !strings.HasSuffix(toLower, "l") {
		mult = 1.0 / 1000.0
	}
	amnt, err := strconv.ParseFloat(strings.Trim(toLower, "ml"), 64)
	if err != nil {
		msg.Return("I couldn't parse the liquid amount! | Usage: .waliquor [size] [amount]")
		return
	}
	amnt = amnt * mult
	literTax := (amnt * 3.7708)
	salesTax := (dollar * 0.205)
	totalCost := literTax + salesTax + dollar
	msg.Return(fmt.Sprintf("Total Cost: $%v (20.5%% Sales Tax: $%v, $3.7708/L Tax: $%v)", totalCost, salesTax, literTax))
}
