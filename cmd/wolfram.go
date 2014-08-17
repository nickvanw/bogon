package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/clbanning/x2j"
)

func init() {
	AddPlugin("WolframAlpha", "(?i)^\\.w(olfram)?(a(lpha)?)?$", MessageHandler(WolframAlpha), false, false)
}

const geturl = "http://api.wolframalpha.com/v2/query"

func WolframAlpha(msg *Message) {
	appid, avail := GetConfig("WolframAlpha")
	if avail != true {
		fmt.Println("Wolfram AppID Not available")
		return
	}
	query := urlencode(strings.Join(msg.Params[1:], " "))
	getURL := fmt.Sprintf("%s?input=%s&appid=%s", geturl, query, appid)
	resp, err := getSite(getURL)
	if err != nil {
		msg.Return("Wolfram returned an error!")
		return
	}
	data, err := x2j.ByteDocToJson(resp)
	if err != nil {
		msg.Return("Wolfram returned Invalid XML!")
		return
	}
	var wresp Wolfram
	json.Unmarshal([]byte(data), &wresp)
	if len(wresp.Queryresult.Pod) > 1 {
		input := strings.Replace(wresp.Queryresult.Pod[0].Subpod.Plaintext, "\n", " ", -1)
		resp := strings.Replace(wresp.Queryresult.Pod[1].Subpod.Plaintext, "\n", " ", -1)
		msg.Return(fmt.Sprintf("%s: %s", input, resp))
	} else {
		msg.Return("I didn't get any data for that!")
	}

}

type Wolfram struct {
	Queryresult struct {
		_Datatypes     string `json:"-datatypes"`
		_Error         string `json:"-error"`
		_Host          string `json:"-host"`
		_Id            string `json:"-id"`
		_Numpods       string `json:"-numpods"`
		_Parsetimedout string `json:"-parsetimedout"`
		_Parsetiming   string `json:"-parsetiming"`
		_Recalculate   string `json:"-recalculate"`
		_Related       string `json:"-related"`
		_Server        string `json:"-server"`
		_Success       string `json:"-success"`
		_Timedout      string `json:"-timedout"`
		_Timedoutpods  string `json:"-timedoutpods"`
		_Timing        string `json:"-timing"`
		_Version       string `json:"-version"`
		Assumptions    struct {
			_Count     string `json:"-count"`
			Assumption []struct {
				_Count    string `json:"-count"`
				_Template string `json:"-template"`
				_Type     string `json:"-type"`
				_Word     string `json:"-word"`
				Value     []struct {
					_Desc  string `json:"-desc"`
					_Input string `json:"-input"`
					_Name  string `json:"-name"`
				} `json:"value"`
			} `json:"assumption"`
		} `json:"assumptions"`
		Pod []struct {
			_Error      string `json:"-error"`
			_Id         string `json:"-id"`
			_Numsubpods string `json:"-numsubpods"`
			_Position   string `json:"-position"`
			_Scanner    string `json:"-scanner"`
			_Title      string `json:"-title"`
			Subpod      struct {
				_Title string `json:"-title"`
				Img    struct {
					_Alt    string `json:"-alt"`
					_Height string `json:"-height"`
					_Src    string `json:"-src"`
					_Title  string `json:"-title"`
					_Width  string `json:"-width"`
				} `json:"img"`
				Plaintext string `json:"plaintext"`
			} `json:"subpod"`
		} `json:"pod"`
		Sources struct {
			_Count string `json:"-count"`
			Source struct {
				_Text string `json:"-text"`
				_Url  string `json:"-url"`
			} `json:"source"`
		} `json:"sources"`
	} `json:"queryresult"`
}
