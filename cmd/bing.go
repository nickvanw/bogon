package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/crowdmob/goamz/s3"
)

func init() {
	AddPlugin("Bing", "(?i)^\\.(google|bing)$", MessageHandler(Bing), false, false)
	AddPlugin("BingImage", "(?i)^\\.im(a)?g(e)?$", MessageHandler(Bing), false, false)
}

type BingProcess interface {
	Process(data []byte) (string, error)
	ID() string
}

func Bing(msg *Message) {
	api_token, avail := GetConfig("Bing")
	if avail != true {
		fmt.Println("Bing API Key not Found")
		return
	}
	var BingSearchProcess BingProcess
	switch msg.Name {
	case "Bing":
		BingSearchProcess = BingSearch{}
	case "BingImage":
		BingSearchProcess = BingImage{}
	default:
		fmt.Println("Unknown Bing Call!")
		return
	}
	query := strings.Join(msg.Params[1:], " ")
	bing, err := GetBing(query, BingSearchProcess, api_token)
	if err != nil {
		switch err.Error() {
		case "1":
			msg.Return("I had an error trying to search Bing, sorry!")
			return
		case "0":
			msg.Return("I couldn't find any results, sorry!")
			return
		}
	}
	msg.Return(bing)
}
func GetBing(search string, processor BingProcess, api_token string) (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.datamarket.azure.com/Bing/Search/%s?$format=json&Adult=%%27Off%%27&$top=1&Query=%%27%s%%27", processor.ID(), url.QueryEscape(search))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.New("1")
	}
	req.SetBasicAuth(api_token, api_token)
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("1")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("1")
	}
	data, err := processor.Process(body)
	if err != nil {
		return "", errors.New("0")
	}
	return data, nil
}

type BingImage struct{}
type BingSearch struct{}

func (bi BingImage) ID() string {
	return "Image"
}

func (bi BingSearch) ID() string {
	return "Web"
}
func (bi BingImage) Process(data []byte) (string, error) {
	var resp ImageResponse
	json.Unmarshal(data, &resp)
	if len(resp.D.Results) < 1 {
		return "", errors.New("0")
	}
	image, err := getSite(resp.D.Results[0].MediaUrl)
	if err != nil {
		return "", errors.New("0")
	}
	fileName := GenerateRandom(6)
	uploaded, err := UploadToS3(fileName, image, resp.D.Results[0].ContentType, s3.PublicRead)
	if err != nil {
		return "", errors.New("0")
	} else {
		out := fmt.Sprintf("Image Result: %s", uploaded)
		return out, nil
	}
}

func (bi BingSearch) Process(data []byte) (string, error) {
	var resp BingResponse
	json.Unmarshal(data, &resp)
	if len(resp.D.Results) < 1 {
		return "", errors.New("0")
	}
	out := fmt.Sprintf("%s: %s", resp.D.Results[0].Title, resp.D.Results[0].URL)
	return out, nil
}

type ImageResponse struct {
	D struct {
		Next    string `json:"__next"`
		Results []struct {
			ContentType string `json:"ContentType"`
			DisplayUrl  string `json:"DisplayUrl"`
			FileSize    string `json:"FileSize"`
			Height      string `json:"Height"`
			ID          string `json:"ID"`
			MediaUrl    string `json:"MediaUrl"`
			SourceUrl   string `json:"SourceUrl"`
			Thumbnail   struct {
				ContentType string `json:"ContentType"`
				FileSize    string `json:"FileSize"`
				Height      string `json:"Height"`
				MediaUrl    string `json:"MediaUrl"`
				Width       string `json:"Width"`
				Metadata    struct {
					Type string `json:"type"`
				} `json:"__metadata"`
			} `json:"Thumbnail"`
			Title    string `json:"Title"`
			Width    string `json:"Width"`
			Metadata struct {
				Type string `json:"type"`
				Uri  string `json:"uri"`
			} `json:"__metadata"`
		} `json:"results"`
	} `json:"d"`
}
type BingResponse struct {
	D struct {
		Next    string `json:"__next"`
		Results []struct {
			Description string `json:"Description"`
			DisplayUrl  string `json:"DisplayUrl"`
			ID          string `json:"ID"`
			Title       string `json:"Title"`
			URL         string `json:"Url"`
			Metadata    struct {
				Type string `json:"type"`
				Uri  string `json:"uri"`
			} `json:"__metadata"`
		} `json:"results"`
	} `json:"d"`
}
