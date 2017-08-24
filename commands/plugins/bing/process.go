package bing

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dchest/uniuri"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/bogon/commands/util"
)

type bingProcesser interface {
	sType() string
	sParams() string
	process(data []byte) (string, error)
}

func bingAPIFetch(query, token string, p bingProcesser) (string, error) {
	url := fmt.Sprintf("https://api.cognitive.microsoft.com/bing/v5.0/%s?safeSearch=Off&count=1&q=%s", p.sType(), url.QueryEscape(query))
	client := new(http.Client)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", token)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return p.process(data)
}

type bingSearchProcess struct{}
type bingImageProcess struct{}

func (bingImageProcess) sType() string {
	return "/images/search"
}

func (bingSearchProcess) sType() string {
	return "/search"
}

func (bingImageProcess) sParams() string {
	return ""
}

func (bingSearchProcess) sParams() string {
	return "&responseFilter=Webpages"
}

func (bingSearchProcess) process(data []byte) (string, error) {
	var resp searchResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", err
	}
	if len(resp.WebPages.Value) < 1 {
		return "", errors.New("no results")
	}
	out := fmt.Sprintf("%s: %s", resp.WebPages.Value[0].Name, resp.WebPages.Value[0].URL)
	return out, nil
}

func (bingImageProcess) process(data []byte) (string, error) {
	var api imageResponse
	if err := json.Unmarshal(data, &api); err != nil {
		return "", err
	}
	if len(api.Value) == 0 {
		return "", errors.New("no results")
	}
	img := api.Value[0]
	client := new(http.Client)
	req, err := http.NewRequest("GET", img.ContentURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	name := uniuri.NewLen(10)
	name += fmt.Sprintf(".%s", img.EncodingFormat) // EncodingFormat is jpeg/png/etc.
	bucketName, bucketCfg := config.Get("BING_S3_BUCKET")
	awsRegion, regionCfg := config.Get("BING_S3_REGION")
	if !bucketCfg || !regionCfg {
		return "", errors.New("invalid bucket/region to upload")
	}
	url, err := util.UploadWithEnv(bucketName, awsRegion, name, img.EncodingFormat, resp.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Image Result: %s", url), nil
}

type searchResponse struct {
	WebPages struct {
		Value []struct {
			Name string `json:"name"`
			URL  string `json:"displayUrl"`
		} `json:"value"`
	} `json:"webPages"`
}

type imageResponse struct {
	Value []struct {
		ContentURL     string `json:"contentUrl"`
		EncodingFormat string `json:"encodingFormat"`
	} `json:"value"`
}
