package youtube

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/ixai/iso8601duration"
	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/bogon/commands/util"
)

const (
	maxPerLine  = 2
	timeForm    = "2006-01-02T15:04:05.000Z"
	baseURL     = "https://www.googleapis.com/youtube/v3/"
	videoParts  = "contentDetails,snippet,statistics"
	videoFields = "items(contentDetails(duration),snippet(channelId,publishedAt,title),statistics(viewCount))"
	chanParts   = "snippet"
	chanFields  = "items(snippet(title))"
	ytTitle     = "youtube"
)

// YoutubeCommand will register the handlers necessary to scan for youtube URLs
var YoutubeCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	return ytTitle, nil, handleYT, commands.Options{Raw: true}
}

func handleYT(msg commands.Message, ret commands.MessageFunc) string {
	r := regexp.MustCompile(`(?i)(?:v=|\/)(([\w-]+){11})`)
	vid := 0
	for _, v := range msg.Params {
		match := r.FindAllStringSubmatch(v, -1)
		if match != nil && vid < maxPerLine {
			vid++
			id := match[0][1]
			yt, ex := getYT(id)
			if ex {
				ret(yt)
			}
		}
	}
	return ""
}

func getYT(msg string) (string, bool) {
	if data, is := getVideoInfo(msg); is {
		return formatYT(data), true
	}
	return "", false
}

func getVideoInfo(id string) (*ytVideoResponse, bool) {
	key, avail := config.Get("GOOGLE_API")
	if !avail {
		return nil, false
	}

	videoURL := fmt.Sprintf("%svideos/?id=%s&key=%s&part=%s&fields=%s", baseURL, id, key, videoParts, videoFields)
	vData, err := util.Fetch(videoURL)
	if err != nil {
		return nil, false
	}

	var videoData tyVideo
	err = json.Unmarshal(vData, &videoData)
	if err != nil || len(videoData.Videos) < 1 {
		return nil, false
	}

	chanID := videoData.Videos[0].Snippet.ChannelID
	chanURL := fmt.Sprintf("%schannels/?id=%s&key=%s&part=%s&fields=%s", baseURL, chanID, key, chanParts, chanFields)
	cData, err := util.Fetch(chanURL)
	if err != nil {
		return nil, false
	}

	var chanData ytChannel
	err = json.Unmarshal(cData, &chanData)
	if err != nil || len(chanData.Channels) < 1 {
		return nil, false
	}

	title := videoData.Videos[0].Snippet.Title
	views, _ := strconv.ParseInt(videoData.Videos[0].Statistics.Views, 10, 64)
	uploadTime, _ := time.Parse(timeForm, videoData.Videos[0].Snippet.Published)

	dur, _ := duration.ParseString(videoData.Videos[0].ContentDetails.Duration)
	durTime := dur.ToDuration()

	channel := chanData.Channels[0].Snippet.Title

	response := &ytVideoResponse{
		Title:      title,
		Views:      views,
		Duration:   durTime,
		Channel:    channel,
		UploadTime: uploadTime,
	}

	return response, true
}

func formatYT(yt *ytVideoResponse) string {
	return fmt.Sprintf("%s | Views: %s | Duration: %s | Uploaded By: %s on %s",
		yt.Title, humanize.Comma(yt.Views), yt.Duration, yt.Channel, yt.UploadTime.Format("Jan 2, 2006"))
}

type ytVideoResponse struct {
	Title      string
	Views      int64
	Duration   time.Duration
	Channel    string
	UploadTime time.Time
}

type tyVideo struct {
	Videos []struct {
		Snippet struct {
			Title     string `json:"title"`
			ChannelID string `json:"channelId"`
			Published string `json:"publishedAt"`
		} `json:"snippet"`
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
		Statistics struct {
			Views string `json:"viewCount"`
		} `json:"statistics"`
	} `json:"items"`
}

type ytChannel struct {
	Channels []struct {
		Snippet struct {
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}
