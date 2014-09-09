package cmd

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "strings"
)

func init() {
        AddPlugin("Shorten", "(?i)^\\.short(en)?$", MessageHandler(Shorten), false, false)
}

const apiUrl = "https://www.googleapis.com/urlshortener/v1/url"

func Shorten(msg *Message) {
        jsonOut := fmt.Sprintf("{\"longUrl\": \"%s\"}", strings.Join(msg.Params[1:], " "))
        req, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(jsonOut))
        req.Header.Set("Content-Type", "application/json")
        
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
                msg.Return("Error contacting Google URL Shortener!")
                return
        }
        defer resp.Body.Close()
        
        body, _ := ioutil.ReadAll(resp.Body)
        
        var si ShortenInfo
        json.Unmarshal(body, &si)
        
        short := si.Short
        
        msg.Return(short)
}

type ShortenInfo struct {
        Kind    string  `json:"kind"`
        Short   string  `json:"id"`
        Long    string  `json:"longUrl"`
}
