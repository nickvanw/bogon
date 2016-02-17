package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/util"
)

var ipLookup = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.(ip|geo)$")
	return ipTitle, out, ipAddressLookup, defaultOptions
}

func ipAddressLookup(msg commands.Message, ret commands.MessageFunc) string {
	ip := strings.Join(msg.Params[1:], " ")
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	data, err := util.Fetch(url)
	if err != nil {
		return "sorry, there was a geolookup error.."
	}
	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		return "Sorry, ip-api.com gave me a bad response"
	}
	return fmt.Sprintf("%s: %s, %s, %s, %s, %s with %s",
		response["query"], response["city"], response["regionName"],
		response["country"], response["org"], response["as"], response["isp"])
}
