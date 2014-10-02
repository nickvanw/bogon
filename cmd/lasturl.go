package cmd

import (
        "net/url"
        "strings"
)

func HandleLastUrl(msg *Message) {
    for _, v := range msg.Params {
        u, err := url.Parse(v)
        
        if err == nil {
            url := u.String()
        
            if (strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
                channel, err := msg.State.GetChan(msg.To)
                if err == nil {
                    msg.State.Lock()
                    channel.LastUrl = url
                    msg.State.Unlock()
                    break
                }
            }
        }
    }
}