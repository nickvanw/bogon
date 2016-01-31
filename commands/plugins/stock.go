package plugins

import (
	"encoding/csv"
	"fmt"
	"regexp"
	"strings"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/util"
)

var stockCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.stock$")
	return stockTitle, out, stockLookup, defaultOptions
}

func stockLookup(msg commands.Message, ret commands.MessageFunc) string {
	stock := strings.Join(msg.Params[1:], "%20")
	url := fmt.Sprintf("http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=snl1abcghm6m8", stock)
	webdata, err := util.Fetch(url)
	if err != nil {
		return "Error fetching that stock!"
	}
	reader := strings.NewReader(string(webdata))
	csvReader := csv.NewReader(reader)
	csvData, err := csvReader.ReadAll()
	if err != nil || len(csvData) < 1 {
		return "Error fetching that stock!"
	}
	data := csvData[0]
	if string(data[3]) == "N/A" && string(data[4]) == "N/A" && string(data[6]) == "N/A" {
		return "Invalid Stock!"
	}
	return fmt.Sprintf("%s (%s): Last: $%s | Ask: $%s | Bid: $%s | Change: %s | Day Low: %s | Day High: %s | %s Last 200 days | %s Last 50 days",
		data[1], data[0], data[2], data[3], data[4], data[5], data[6], data[7], data[8], data[9])

}
