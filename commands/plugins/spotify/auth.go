package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/nickvanw/bogon/commands/config"
)

type spotifyAuth struct {
	token     string
	expiresAt time.Time
	updatedAt time.Time
	sync.RWMutex
}

var tokenURL = "https://accounts.spotify.com/api/token"

func (s *spotifyAuth) fetch(url string) ([]byte, error) {
	if err := s.refresh(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	s.RLock()
	req.Header.Add("Authorization", "Bearer "+s.token)
	s.RUnlock()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *spotifyAuth) refresh() error {
	s.RLock()
	needsUpdate := s.updatedAt.IsZero() || time.Now().After(s.expiresAt)
	s.RUnlock()
	if needsUpdate {
		return s.fetchSecret()
	}
	return nil
}

func (s *spotifyAuth) fetchSecret() error {
	creds, ok := s.spotifySecret()
	if !ok {
		return errors.New("No Spotify Credentials Configured")
	}
	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", creds))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	var authData authResponse
	if err := json.NewDecoder(resp.Body).Decode(&authData); err != nil {
		return err
	}
	s.Lock()
	s.expiresAt = time.Now().Add(time.Duration(authData.ExpiresIn) * time.Second)
	s.updatedAt = time.Now()
	s.token = authData.AccessToken
	s.Unlock()
	return nil
}

func (s *spotifyAuth) spotifySecret() (string, bool) {
	clientID, cOK := config.Get("SPOTIFY_CLIENT_ID")
	clientSecret, sOK := config.Get("SPOTIFY_CLIENT_SECRET")

	if !cOK || !sOK {
		return "", false
	}

	return base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret)), true
}

type authResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}
