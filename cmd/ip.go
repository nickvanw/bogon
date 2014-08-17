package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
)

func init() {
	AddPlugin("IPLookup", "(?i)^\\.(ip|geo)$", MessageHandler(IPLookup), false, false)
}

func IPLookup(msg *Message) {
	ip := strings.Join(msg.Params[1:], " ")
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	data, err := getSite(url)
	if err != nil {
		msg.Return("Error!")
		return
	}
	var response map[string]string
	json.Unmarshal(data, &response)
	out := fmt.Sprintf("%v: %v, %v, %v, %v, %v with %v",
		response["query"], response["city"], response["regionName"],
		response["country"], response["org"], response["as"], response["isp"])
	msg.Return(out)
}
