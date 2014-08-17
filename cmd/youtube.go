package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
)

const timeForm = "2006-01-02T15:04:05.000Z"

func GetVideoInfo(id string) (YouTubeVideoResponse, bool) {
	url := fmt.Sprintf("http://gdata.youtube.com/feeds/api/videos/%s?alt=json", id)
	data, _ := getSite(url)
	var videoData YouTubeVideo
	er := json.Unmarshal(data, &videoData)
	if er != nil {
		return YouTubeVideoResponse{}, false
	}
	title := videoData.Entry.Title.T
	author := videoData.Entry.Author[0].Name.T
	views, _ := strconv.Atoi(videoData.Entry.Yt_Statistics.ViewCount)
	views64 := int64(views)
	length := videoData.Entry.Media_Group.Yt_Duration.Seconds
	durationString := fmt.Sprintf("%ss", length)
	duration, _ := time.ParseDuration(durationString)
	uploadDate := videoData.Entry.Published.T
	uploadTime, _ := time.Parse(timeForm, uploadDate)
	response := YouTubeVideoResponse{
		Title:    title,
		Author:   author,
		Views:    views64,
		Uploaded: uploadTime,
		Duration: duration,
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
	msg := fmt.Sprintf("%s | Views: %s | Duration: %s | Uploaded By: %s on %s", yt.Title, humanize.Comma(yt.Views), yt.Duration, yt.Author, yt.Uploaded.Format("Jan 2, 2006"))
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
	Title    string
	Author   string
	Views    int64
	Uploaded time.Time
	Duration time.Duration
}
type YouTubeVideo struct {
	Encoding string `json:"encoding"`
	Entry    struct {
		App_Control struct {
			Yt_State struct {
				T          string `json:"$t"`
				Name       string `json:"name"`
				ReasonCode string `json:"reasonCode"`
			} `json:"yt$state"`
		} `json:"app$control"`
		Author []struct {
			Name struct {
				T string `json:"$t"`
			} `json:"name"`
			Uri struct {
				T string `json:"$t"`
			} `json:"uri"`
		} `json:"author"`
		Category []struct {
			Scheme string `json:"scheme"`
			Term   string `json:"term"`
		} `json:"category"`
		Content struct {
			T    string `json:"$t"`
			Type string `json:"type"`
		} `json:"content"`
		Gd_Comments struct {
			Gd_FeedLink struct {
				CountHint float64 `json:"countHint"`
				Href      string  `json:"href"`
				Rel       string  `json:"rel"`
			} `json:"gd$feedLink"`
		} `json:"gd$comments"`
		Gd_Rating struct {
			Average   float64 `json:"average"`
			Max       float64 `json:"max"`
			Min       float64 `json:"min"`
			NumRaters float64 `json:"numRaters"`
			Rel       string  `json:"rel"`
		} `json:"gd$rating"`
		ID struct {
			T string `json:"$t"`
		} `json:"id"`
		Link []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
			Type string `json:"type"`
		} `json:"link"`
		Media_Group struct {
			Media_Category []struct {
				T      string `json:"$t"`
				Label  string `json:"label"`
				Scheme string `json:"scheme"`
			} `json:"media$category"`
			Media_Content []struct {
				Duration   float64 `json:"duration"`
				Expression string  `json:"expression"`
				IsDefault  string  `json:"isDefault"`
				Medium     string  `json:"medium"`
				Type       string  `json:"type"`
				URL        string  `json:"url"`
				Yt_Format  float64 `json:"yt$format"`
			} `json:"media$content"`
			Media_Description struct {
				T    string `json:"$t"`
				Type string `json:"type"`
			} `json:"media$description"`
			Media_Keywords struct{} `json:"media$keywords"`
			Media_Player   []struct {
				URL string `json:"url"`
			} `json:"media$player"`
			Media_Thumbnail []struct {
				Height float64 `json:"height"`
				Time   string  `json:"time"`
				URL    string  `json:"url"`
				Width  float64 `json:"width"`
			} `json:"media$thumbnail"`
			Media_Title struct {
				T    string `json:"$t"`
				Type string `json:"type"`
			} `json:"media$title"`
			Yt_Duration struct {
				Seconds string `json:"seconds"`
			} `json:"yt$duration"`
		} `json:"media$group"`
		Published struct {
			T string `json:"$t"`
		} `json:"published"`
		Title struct {
			T    string `json:"$t"`
			Type string `json:"type"`
		} `json:"title"`
		Updated struct {
			_T string `json:"$t"`
		} `json:"updated"`
		Xmlns         string   `json:"xmlns"`
		Xmlns_App     string   `json:"xmlns$app"`
		Xmlns_Gd      string   `json:"xmlns$gd"`
		Xmlns_Media   string   `json:"xmlns$media"`
		Xmlns_Yt      string   `json:"xmlns$yt"`
		Yt_Hd         struct{} `json:"yt$hd"`
		Yt_Statistics struct {
			FavoriteCount string `json:"favoriteCount"`
			ViewCount     string `json:"viewCount"`
		} `json:"yt$statistics"`
	} `json:"entry"`
	Version string `json:"version"`
}
