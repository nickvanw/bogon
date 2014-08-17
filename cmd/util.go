package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getSite(url string) ([]byte, error) {
	resp, err := http.Get(url)
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

func bold(s string) string {
	return fmt.Sprintf("\x02%s\x02", s)
}
