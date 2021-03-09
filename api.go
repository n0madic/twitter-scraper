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
			Legacy struct {
				CreatedAt   string `json:"created_at"`
				Description string `json:"description"`
				Entities    struct {
					URL struct {
						Urls []struct {
							ExpandedURL string `json:"expanded_url"`
						} `json:"urls"`
					} `json:"url"`
				} `json:"entities"`
				FavouritesCount      int      `json:"favourites_count"`
				FollowersCount       int      `json:"followers_count"`
				FriendsCount         int      `json:"friends_count"`
				IDStr                string   `json:"id_str"`
				ListedCount          int      `json:"listed_count"`
				Name                 string   `json:"name"`
				Location             string   `json:"location"`
				PinnedTweetIdsStr    []string `json:"pinned_tweet_ids_str"`
				ProfileBannerURL     string   `json:"profile_banner_url"`
				ProfileImageURLHTTPS string   `json:"profile_image_url_https"`
				Protected            bool     `json:"protected"`
				ScreenName           string   `json:"screen_name"`
				StatusesCount        int      `json:"statuses_count"`
				Verified             bool     `json:"verified"`
			} `json:"legacy"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
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

	// private profiles return forbidden, but also data
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
		return fmt.Errorf("response status %s", resp.Status)
	}

	if resp.Header.Get("X-Rate-Limit-Remaining") == "0" {
		s.guestToken = ""
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
