package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Fetch will attempt to call and fetch the specified URL,
// returning the body in a byte array
func Fetch(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// URLEncode prepares a string to be sent over HTTP
func URLEncode(s string) (result string) {
	for _, c := range s {
		if c <= 0x7f { // single byte
			result += fmt.Sprintf("%%%X", c)
		} else if c > 0x1fffff { // quaternary byte
			result += fmt.Sprintf("%%%X%%%X%%%X%%%X",
				0xf0+((c&0x1c0000)>>18),
				0x80+((c&0x3f000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else if c > 0x7ff { // triple byte
			result += fmt.Sprintf("%%%X%%%X%%%X",
				0xe0+((c&0xf000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else { // double byte
			result += fmt.Sprintf("%%%X%%%X",
				0xc0+((c&0x7c0)>>6),
				0x80+(c&0x3f),
			)
		}
	}

	return result
}

// Bold returns the string wrapped in the bold control codes for IRC
func Bold(s string) string {
	return fmt.Sprintf("\x02%s\x02", s)
}

// StripNewLines strips the newlines from the end of strings
func StripNewLines(in string) string {
	rn := strings.Replace(in, "\r\n", " ", -1)
	n := strings.Replace(rn, "\n", " ", -1)
	return strings.Replace(n, "\r", " ", -1)
}
