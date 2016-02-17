package bing

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"

	"github.com/dchest/uniuri"
	"github.com/nickvanw/bogon/commands/util"
)

const bucketName = "img"

type bingProcesser interface {
	sType() string
	process(data []byte) (string, error)
}

func bingAPIFetch(query, token string, p bingProcesser) (string, error) {
	url := fmt.Sprintf("https://api.datamarket.azure.com/Bing/Search/%s?$format=json&Adult=%%27Off%%27&$top=1&Query=%%27%s%%27",
		p.sType(), url.QueryEscape(query))
	client := new(http.Client)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(token, token)
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
	return "Image"
}

func (bingSearchProcess) sType() string {
	return "Web"
}

func (bingSearchProcess) process(data []byte) (string, error) {
	var resp searchResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", err
	}
	if len(resp.D.Results) < 1 {
		return "", errors.New("no results")
	}
	out := fmt.Sprintf("%s: %s", resp.D.Results[0].Title, resp.D.Results[0].URL)
	return out, nil
}

func (bingImageProcess) process(data []byte) (string, error) {
	var api imageResponse
	if err := json.Unmarshal(data, &api); err != nil {
		return "", err
	}
	if len(api.D.Results) == 0 {
		return "", errors.New("no results")
	}
	img := api.D.Results[0]
	client := new(http.Client)
	req, err := http.NewRequest("GET", img.MediaURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	name := uniuri.NewLen(10)
	extention, err := mime.ExtensionsByType(img.ContentType)
	if err != nil || len(extention) == 0 {
		name += ".png" // browser should pick up on it anyway
	} else {
		name += extention[0]
	}
	url, err := util.UploadWithEnv(bucketName, name, img.ContentType, resp.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Image Result: %s", url), nil
}

type searchResponse struct {
	D struct {
		Results []struct {
			Title string `json:"Title"`
			URL   string `json:"Url"`
		} `json:"results"`
	} `json:"d"`
}

type imageResponse struct {
	D struct {
		Results []struct {
			ContentType string `json:"ContentType"`
			MediaURL    string `json:"MediaUrl"`
		} `json:"results"`
	} `json:"d"`
}
