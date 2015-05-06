package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/ixai/iso8601duration"
)

const timeForm = "2006-01-02T15:04:05.000Z"

const baseUrl = "https://www.googleapis.com/youtube/v3/"
const videoParts = "contentDetails,snippet,statistics"
const videoFields = "items(contentDetails(duration),snippet(channelId,publishedAt,title),statistics(viewCount))"
const chanParts = "snippet"
const chanFields = "items(snippet(title))"

func GetVideoInfo(id string) (YouTubeVideoResponse, bool) {
	key, avail := GetConfig("Google")
	if !avail {
		return YouTubeVideoResponse{}, false
	}

	videoUrl := fmt.Sprintf("%svideos/?id=%s&key=%s&part=%s&fields=%s", baseUrl, id, key, videoParts, videoFields)
	vData, _ := getSite(videoUrl)

	var videoData YouTubeVideo
	err := json.Unmarshal(vData, &videoData)
	if err != nil || len(videoData.Videos) < 1 {
		return YouTubeVideoResponse{}, false
	}

	chanId := videoData.Videos[0].Snippet.ChannelId
	chanUrl := fmt.Sprintf("%schannels/?id=%s&key=%s&part=%s&fields=%s", baseUrl, chanId, key, chanParts, chanFields)
	cData, _ := getSite(chanUrl)

	var chanData YouTubeChannel
	err = json.Unmarshal(cData, &chanData)
	if err != nil || len(chanData.Channels) < 1 {
		return YouTubeVideoResponse{}, false
	}

	title := videoData.Videos[0].Snippet.Title
	views, _ := strconv.ParseInt(videoData.Videos[0].Statistics.Views, 10, 64)
	uploadTime, _ := time.Parse(timeForm, videoData.Videos[0].Snippet.Published)

	dur, _ := duration.ParseString(videoData.Videos[0].ContentDetails.Duration)
	durTime := dur.ToDuration()

	channel := chanData.Channels[0].Snippet.Title

	response := YouTubeVideoResponse{
		Title:      title,
		Views:      views,
		Duration:   durTime,
		Channel:    channel,
		UploadTime: uploadTime,
	}

	return response, true
}

func HandleYoutube(msg []string, out *Message) {
	r := regexp.MustCompile(`(?i)[v=|\/]([\w-]+)(&.+)?$`)
	vid := 0
	for _, v := range msg {
		match := r.FindAllStringSubmatch(v, -1)
		if match != nil && vid < 2 {
			vid = vid + 1
			id := match[0][1]
			yt, ex := GetYoutube(id)
			if ex {
				out.Return(yt)
			}
		}
	}
}

func FormatYouTube(yt YouTubeVideoResponse) string {
	msg := fmt.Sprintf("%s | Views: %s | Duration: %s | Uploaded By: %s on %s", yt.Title, humanize.Comma(yt.Views), yt.Duration, yt.Channel, yt.UploadTime.Format("Jan 2, 2006"))
	return msg
}

func GetYoutube(msg string) (string, bool) {
	data, is := GetVideoInfo(msg)
	if is == true {
		return FormatYouTube(data), true
	}
	return "", false
}

type YouTubeVideoResponse struct {
	Title      string
	Views      int64
	Duration   time.Duration
	Channel    string
	UploadTime time.Time
}

type YouTubeVideo struct {
	Videos []struct {
		Snippet struct {
			Title     string `json:"title"`
			ChannelId string `json:"channelId"`
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

type YouTubeChannel struct {
	Channels []struct {
		Snippet struct {
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}
