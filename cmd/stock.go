package cmd

import (
	"encoding/csv"
	"fmt"
	"strings"
)

func init() {
	AddPlugin("Stock", "(?i)^\\.stock$", MessageHandler(Stock), false, false)
}

func Stock(msg *Message) {
	stock := strings.Join(msg.Params[1:], "%20")
	url := fmt.Sprintf("http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=snl1abcghm6m8", stock)
	webdata, err := getSite(url)
	if err != nil {
		msg.Return("Error fetching that stock!")
		return
	}
	reader := strings.NewReader(string(webdata))
	csvReader := csv.NewReader(reader)
	csvData, err := csvReader.ReadAll()
	if err != nil || len(csvData) < 1 {
		msg.Return("Error fetching that stock!")
		return
	}
	data := csvData[0]
	if string(data[3]) == "N/A" && string(data[4]) == "N/A" && string(data[6]) == "N/A" {
		msg.Return("Invalid Stock!")
		return
	}
	retval := fmt.Sprintf("%s (%s): Last: $%s | Ask: $%s | Bid: $%s | Change: %s | Day Low: %s | Day High: %s | %s Last 200 days | %s Last 50 days",
		data[1], data[0], data[2], data[3], data[4], data[5], data[6], data[7], data[8], data[9])
	msg.Return(retval)

}
