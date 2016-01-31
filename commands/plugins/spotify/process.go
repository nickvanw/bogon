package spotify

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nickvanw/bogon/commands/util"
)

func artistLookup(key string) string {
	fullURL := fmt.Sprintf("%s/%s/%s", lookupURL, "artists", key)
	data, err := util.Fetch(fullURL)
	if err != nil {
		return ""
	}
	var ret spotifyArtist
	if err := json.Unmarshal(data, &ret); err != nil {
		return ""
	}
	return fmt.Sprintf("[Artist]: %s with %d followers and %d popularity", ret.Name, ret.Followers.Total, ret.Popularity)
}

type spotifyArtist struct {
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
}

func albumLookup(key string) string {
	fullURL := fmt.Sprintf("%s/%s/%s", lookupURL, "albums", key)
	data, err := util.Fetch(fullURL)
	if err != nil {
		return ""
	}
	var ret spotifyAlbum
	if err := json.Unmarshal(data, &ret); err != nil {
		return ""
	}
	var artist string
	if len(ret.Artists) > 1 {
		artist = "Multiple Artists"
	} else {
		artist = ret.Artists[0].Name
	}
	return fmt.Sprintf("[Album]: %s by %s released in %s", ret.Name, artist, ret.ReleaseDate)
}

type spotifyAlbum struct {
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
}

func trackLookup(key string) string {
	fullURL := fmt.Sprintf("%s/%s/%s", lookupURL, "tracks", key)
	data, err := util.Fetch(fullURL)
	if err != nil {
		return ""
	}
	var ret spotifyTrack
	if err := json.Unmarshal(data, &ret); err != nil {
		return ""
	}

	duration := time.Duration(ret.DurationMs) * time.Millisecond
	var artist string
	if len(ret.Artists) > 1 {
		artist = "Multiple Artists"
	} else {
		artist = ret.Artists[0].Name
	}
	return fmt.Sprintf("[Song]: %s - %s (%s) on the album %s", ret.Name, artist, duration.String(), ret.Album.Name)
}

type spotifyTrack struct {
	Album struct {
		Name string `json:"name"`
	} `json:"album"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	DurationMs int    `json:"duration_ms"`
	Name       string `json:"name"`
}
