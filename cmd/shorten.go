package cmd

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "strings"
        "time"
)

func init() {
        AddPlugin("Shorten", "(?i)^\\.short(en)?$", MessageHandler(Shorten), false, false)
}

const apiUrl = "https://www.googleapis.com/urlshortener/v1/url"

func Shorten(msg *Message) {
        var longUrl string
        if len(msg.Params) < 2 {
            channel, err := msg.State.GetChan(msg.To)
            if err != nil {
                msg.Return("Couldn't find the channel I'm in!")
                return
            }
            longUrl = channel.LastUrl
        } else {
            longUrl = strings.Join(msg.Params[1:], " ")
        }

        jsonOut := fmt.Sprintf("{\"longUrl\": \"%s\"}", longUrl)
        req, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(jsonOut))
        req.Header.Set("Content-Type", "application/json")
        
        client := &http.Client{
                Timeout: time.Second * 3,
        }
        resp, err := client.Do(req)
        if err != nil {
                msg.Return("Error contacting Google URL Shortener!")
                return
        }
        defer resp.Body.Close()
        
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                msg.Return("Error reading data from Google URL Shortener!")
                return
        }
        
        var si ShortenInfo
        json.Unmarshal(body, &si)
        
        msg.Return(fmt.Sprintf("Shortened URL: %s", si.Short))
}

type ShortenInfo struct {
        Kind    string  `json:"kind"`
        Short   string  `json:"id"`
        Long    string  `json:"longUrl"`
}
