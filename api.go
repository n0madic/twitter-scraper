package twitterscraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const bearerToken string = "AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"

type user struct {
	Data struct {
		User struct {
			RestID string `json:"rest_id"`
		} `json:"user"`
	} `json:"data"`
}

// Global cache for user IDs
var cacheIDs sync.Map

// RequestAPI get JSON from frontend API and decodes it
func (s *Scraper) RequestAPI(req *http.Request, target interface{}) error {
	if s.guestToken == "" || s.guestCreatedAt.Before(time.Now().Add(-time.Hour*3)) {
		err := s.GetGuestToken()
		if err != nil {
			return err
		}
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("X-Guest-Token", s.guestToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// GetGuestToken from Twitter API
func (s *Scraper) GetGuestToken() error {
	req, err := http.NewRequest("POST", "https://api.twitter.com/1.1/guest/activate.json", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var jsn map[string]interface{}
	if err := json.Unmarshal(body, &jsn); err != nil {
		return err
	}
	var ok bool
	if s.guestToken, ok = jsn["guest_token"].(string); !ok {
		return fmt.Errorf("guest_token not found")
	}
	s.guestCreatedAt = time.Now()

	return nil
}

// GetUserIDByScreenName from API
func (s *Scraper) GetUserIDByScreenName(screenName string) (string, error) {
	id, ok := cacheIDs.Load(screenName)
	if ok {
		return id.(string), nil
	}

	var jsn user
	req, err := http.NewRequest("GET", "https://api.twitter.com/graphql/4S2ihIKfF3xhp-ENxvUAfQ/UserByScreenName?variables=%7B%22screen_name%22%3A%22"+screenName+"%22%2C%22withHighlightedLabel%22%3Atrue%7D", nil)
	if err != nil {
		return "", err
	}

	err = s.RequestAPI(req, &jsn)
	if err != nil {
		return "", err
	}

	if jsn.Data.User.RestID == "" {
		return "", fmt.Errorf("rest_id not found")
	}

	cacheIDs.Store(screenName, jsn.Data.User.RestID)

	return jsn.Data.User.RestID, nil
}
