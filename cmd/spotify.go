package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func HandleSpotify(msg *Message) {
	spouri := regexp.MustCompile(`(?i)(spotify:(?:(?:artist|album|track|user:[^:]+:playlist):[a-zA-Z0-9]+|user:[^:]+|search:(?:[-\w$\.+!*'(),]+|%[a-fA-F0-9]{2})+))`)
	spohttp := regexp.MustCompile(`(?i)http(?:s)?:\/\/(?:open|play)\.spotify\.com\/(track|album|artist)\/([\w]+)`)
	var matches []string
	for _, v := range msg.Params {
		urimatch := spouri.FindAllStringSubmatch(v, -1)
		spohttp := spohttp.FindAllStringSubmatch(v, -1)
		if len(urimatch) > 0 {
			matches = append(matches, urimatch[0][0])
		}
		if len(spohttp) > 0 {
			matches = append(matches, fmt.Sprintf("spotify:%s:%s", spohttp[0][1], spohttp[0][2]))
		}
	}
	if len(matches) > 0 {
		ProcessSpotify(matches, msg)
	}
}

func ProcessSpotify(uri []string, msg *Message) {
	d := 0
	for _, v := range uri {
		if d > 1 {
			return
		}
		d += 1
		uriInfo := strings.Split(v, ":")
		baseURI := "http://ws.spotify.com/lookup/1/.json?uri="
		switch strings.ToLower(uriInfo[1]) {
		case "artist":
			var spotifyData SpotifyArtist
			SpotifyReturn(spotifyData.Process(fmt.Sprintf("%s%s&extras=album", baseURI, v)), msg)
		case "album":
			var spotifyData SpotifyAlbum
			SpotifyReturn(spotifyData.Process(fmt.Sprintf("%s%s", baseURI, v)), msg)
		case "track":
			var spotifyData SpotifyTrack
			SpotifyReturn(spotifyData.Process(fmt.Sprintf("%s%s", baseURI, v)), msg)

		}
	}
}

func SpotifyReturn(data string, msg *Message) {
	msg.Return(data)
}

func (s *SpotifyArtist) Process(uri string) string {
	data, _ := getSite(uri)
	json.Unmarshal(data, &s)
	outData := fmt.Sprintf("(Artist) %s with %v Albums", s.Artist.Name, len(s.Artist.Albums))
	return outData
}

func (s *SpotifyTrack) Process(uri string) string {
	data, _ := getSite(uri)
	json.Unmarshal(data, &s)
	durationString := fmt.Sprintf("%vs", s.Track.Length)
	duration, _ := time.ParseDuration(durationString)
	pop, _ := strconv.ParseFloat(s.Track.Popularity, 64)
	popString := pop * 100
	name := "NaN"
	if len(s.Track.Artists) > 0 {
		name = s.Track.Artists[0].Name
	}
	songData := fmt.Sprintf("(Song): %s - %s (%s)", s.Track.Name, name, duration)
	trackInfo := fmt.Sprintf("Track %s with popularity %v", s.Track.Track_Number, popString)
	albumInfo := fmt.Sprintf("on the album %s, Released in %s", s.Track.Album.Name, s.Track.Album.Released)
	return fmt.Sprintf("%s | %s %s", songData, trackInfo, albumInfo)
}

func (s *SpotifyAlbum) Process(uri string) string {
	data, _ := getSite(uri)
	json.Unmarshal(data, &s)
	outString := fmt.Sprintf("(Album) %s by %s released in %s", s.Album.Name, s.Album.Artist, s.Album.Released)
	return outString
}

type SpotifyTrack struct {
	Info struct {
		Type string `json:"type"`
	} `json:"info"`
	Track struct {
		Album struct {
			Href     string `json:"href"`
			Name     string `json:"name"`
			Released string `json:"released"`
		} `json:"album"`
		Artists []struct {
			Href string `json:"href"`
			Name string `json:"name"`
		} `json:"artists"`
		Availability struct {
			Territories string `json:"territories"`
		} `json:"availability"`
		Available    bool `json:"available"`
		External_Ids []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"external-ids"`
		Href         string  `json:"href"`
		Length       float64 `json:"length"`
		Name         string  `json:"name"`
		Popularity   string  `json:"popularity"`
		Track_Number string  `json:"track-number"`
	} `json:"track"`
}

type SpotifyArtist struct {
	Artist struct {
		Albums []struct {
			Album struct {
				Artist       string `json:"artist"`
				Artist_Id    string `json:"artist-id"`
				Availability struct {
					Territories string `json:"territories"`
				} `json:"availability"`
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"album"`
			Info struct {
				Type string `json:"type"`
			} `json:"info"`
		} `json:"albums"`
		Href string `json:"href"`
		Name string `json:"name"`
	} `json:"artist"`
	Info struct {
		Type string `json:"type"`
	} `json:"info"`
}

type SpotifyAlbum struct {
	Album struct {
		Artist       string `json:"artist"`
		Artist_Id    string `json:"artist-id"`
		Availability struct {
			Territories string `json:"territories"`
		} `json:"availability"`
		External_Ids []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"external-ids"`
		Href     string `json:"href"`
		Name     string `json:"name"`
		Released string `json:"released"`
	} `json:"album"`
	Info struct {
		Type string `json:"type"`
	} `json:"info"`
}
