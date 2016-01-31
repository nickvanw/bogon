package plugins

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/nickvanw/bogon/commands"
)

var waLiquor = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.wal(iquor)?$")
	return waliquorTitle, out, waLiquorCalc, defaultOptions
}

func waLiquorCalc(msg commands.Message, ret commands.MessageFunc) string {
	if len(msg.Params) < 2 {
		return "Usage: .waliquor [size] [amount]"
	}
	dollar, err := strconv.ParseFloat(strings.Replace(msg.Params[len(msg.Params)-1], "$", "", -1), 64)
	if err != nil {
		return "I couldn't parse the currency amount | Usage: .waliquor [size] [amount]"
	}
	toLower := strings.ToLower(msg.Params[1])
	mult := 1.0
	if strings.HasSuffix(toLower, "ml") || !strings.HasSuffix(toLower, "l") {
		mult = 1.0 / 1000.0
	}
	amnt, err := strconv.ParseFloat(strings.Trim(toLower, "ml"), 64)
	if err != nil {
		return "I couldn't parse the liquid amount! | Usage: .waliquor [size] [amount]"
	}
	amnt = amnt * mult
	literTax := (amnt * 3.7708)
	salesTax := (dollar * 0.205)
	totalCost := literTax + salesTax + dollar
	return fmt.Sprintf("Total Cost: $%v (20.5%% Sales Tax: $%v, $3.7708/L Tax: $%v)", totalCost, salesTax, literTax)
}
