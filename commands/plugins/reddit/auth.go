package reddit

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

type redditAuth struct {
	token     string
	expiresAt time.Time
	updatedAt time.Time
	sync.RWMutex
}

var tokenURL = "https://www.reddit.com/api/v1/access_token"

func (r *redditAuth) fetch(url string) ([]byte, error) {
	if err := r.refresh(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	r.RLock()
	req.Header.Add("Authorization", "Bearer "+r.token)
	req.Header.Add("User-Agent", "bogon IRC bot (https://github.com/nickvanw/bogon)")
	r.RUnlock()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Non 2xx status code returned: %d", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *redditAuth) refresh() error {
	r.RLock()
	needsUpdate := r.updatedAt.IsZero() || time.Now().After(r.expiresAt)
	r.RUnlock()
	if needsUpdate {
		return r.fetchSecret()
	}
	return nil
}

func (r *redditAuth) fetchSecret() error {
	creds, ok := r.redditSecret()
	if !ok {
		return errors.New("No Reddit Credentials Configured")
	}
	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", creds))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "bogon IRC bot (https://github.com/nickvanw/bogon)")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	var ret authResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return err
	}

	r.Lock()
	r.expiresAt = time.Now().Add(time.Duration(ret.ExpiresIn) * time.Second)
	r.updatedAt = time.Now()
	r.token = ret.AccessToken
	r.Unlock()
	return nil
}

func (r *redditAuth) redditSecret() (string, bool) {
	clientID, cOK := config.Get("REDDIT_CLIENT_ID")
	clientSecret, sOK := config.Get("REDDIT_CLIENT_SECRET")

	if !cOK || !sOK {
		return "", false
	}

	return base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret)), true
}

type authResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}
