package plugins

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/bogon/commands/util"
)

const geturl = "http://api.wolframalpha.com/v2/query"

var wolframCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.w(olfram)?(a(lpha)?)?$")
	return wolframTitle, out, wolframLookup, defaultOptions
}

func wolframLookup(msg commands.Message, ret commands.MessageFunc) string {
	appid, avail := config.Get("WOLFRAM_API")
	if avail != true {
		return ""
	}
	query := util.URLEncode(strings.Join(msg.Params[1:], " "))
	getURL := fmt.Sprintf("%s?input=%s&appid=%s", geturl, query, appid)
	resp, err := util.Fetch(getURL)
	if err != nil {
		return "Wolfram returned an error!"
	}
	var d result
	if err := xml.Unmarshal(resp, &d); err != nil {
		fmt.Println(err)
		return "Wolfram did not return valid data!"
	}
	if len(d.Pods) < 2 && len(d.Pods[0].Subpods) > 0 && len(d.Pods[1].Subpods) > 0 {
		return "Wolfram did not return valid data!"
	}

	return fmt.Sprintf("wa: q: %s, a: %s", d.Pods[0].Subpods[0].Plaintext, removeNewLineUntil(d.Pods[1].Subpods[0].Plaintext, 200))
}

func removeNewLineUntil(s string, c int) string {
	d := strings.Split(s, "\n")
	if len(d) == 1 && len(s) < c {
		return s
	}
	var out string
	for _, v := range d {
		if len(out)+len(v) > c-3 {
			return out + "..."
		}
		out += " " + v
	}
	return out
}

type result struct {
	XMLName xml.Name `xml:"queryresult"`
	Pods    []pod    `xml:"pod"`
}

type pod struct {
	XMLName xml.Name `xml:"pod"`
	Subpods []subpod `xml:"subpod"`
}

type subpod struct {
	XMLName   xml.Name `xml:"subpod"`
	Plaintext string   `xml:"plaintext"`
}
