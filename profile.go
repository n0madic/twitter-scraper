package twitterscraper

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Profile of twitter user.
type Profile struct {
	Avatar         string
	Banner         string
	Biography      string
	Birthday       string
	FollowersCount int
	FollowingCount int
	IsPrivate      bool
	IsVerified     bool
	Joined         *time.Time
	LikesCount     int
	Location       string
	Name           string
	TweetsCount    int
	URL            string
	UserID         string
	Username       string
	Website        string
}

// GetProfile return parsed user profile.
func GetProfile(username string) (Profile, error) {
	url := "https://mobile.twitter.com/" + username

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Profile{}, err
	}
	req.Header.Set("Accept-Language", "en-US")

	resp, err := http.DefaultClient.Do(req)
	if resp == nil {
		return Profile{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Profile{}, fmt.Errorf("response status: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Profile{}, err
	}

	// parse join date text
	screenName := doc.Find(".screen-name").First().Text()

	// check is username valid
	if screenName == "" {
		return Profile{}, fmt.Errorf("either @%s does not exist or is private", username)
	}

	return Profile{
		Avatar:         doc.Find("td.avatar > img").First().AttrOr("src", ""),
		Biography:      strings.TrimSpace(doc.Find(".bio").First().Text()),
		FollowersCount: parseCount(doc.Find("table.profile-stats > tbody > tr > td:nth-child(3) > a > div.statnum").First().Text()),
		FollowingCount: parseCount(doc.Find("table.profile-stats > tbody > tr > td:nth-child(2) > a > div.statnum").First().Text()),
		IsPrivate:      strings.Contains(doc.Find("div.fullname > a.badge > img").First().AttrOr("src", ""), "protected"),
		IsVerified:     strings.Contains(doc.Find("div.fullname > a.badge > img").First().AttrOr("src", ""), "verified"),
		Location:       strings.TrimSpace(doc.Find(".location").First().Text()),
		Name:           strings.TrimSpace(doc.Find(".fullname").First().Text()),
		TweetsCount:    parseCount(doc.Find("table.profile-stats > tbody > tr > td:nth-child(1) > div.statnum").First().Text()),
		URL:            "https://twitter.com/" + screenName,
		Username:       screenName,
		Website:        strings.TrimSpace(doc.Find("div.url > div > a").First().AttrOr("data-url", "")),
	}, nil
}

func parseCount(str string) (i int) {
	i, _ = strconv.Atoi(strings.Replace(str, ",", "", -1))
	return
}
