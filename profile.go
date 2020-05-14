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
	Biography      string
	Birthday       string
	FollowersCount int
	FollowingCount int
	Joined         *time.Time
	LikesCount     int
	Location       string
	Name           string
	TweetsCount    int
	URL            string
	Username       string
	Website        string
}

// GetProfile return parsed user profile.
func GetProfile(username string) (Profile, error) {
	url := "https://twitter.com/" + username

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

	// parse location, also check is username valid
	location := strings.TrimSpace(doc.Find(".ProfileHeaderCard-locationText.u-dir").First().Text())
	if location == "" {
		return Profile{}, fmt.Errorf("either @%s does not exist or is private", username)
	}

	// parse join date text
	joined, _ := time.Parse("3:4 PM - 2 Jan 2006", doc.Find(".ProfileHeaderCard-joinDateText.u-dir").First().AttrOr("title", ""))

	return Profile{
		Avatar:         doc.Find(".ProfileAvatar-image").First().AttrOr("src", ""),
		Biography:      doc.Find(".ProfileHeaderCard-bio.u-dir").First().Text(),
		Birthday:       strings.ReplaceAll(strings.TrimSpace(doc.Find(".ProfileHeaderCard-birthdateText.u-dir").First().Text()), "Born ", ""),
		FollowersCount: parseCount(doc.Find(".ProfileNav-item--followers > a > span.ProfileNav-value").First()),
		FollowingCount: parseCount(doc.Find(".ProfileNav-item--following > a > span.ProfileNav-value").First()),
		Joined:         &joined,
		LikesCount:     parseCount(doc.Find(".ProfileNav-item--favorites > a > span.ProfileNav-value").First()),
		Location:       location,
		Name:           doc.Find(".ProfileHeaderCard-nameLink").First().Text(),
		TweetsCount:    parseCount(doc.Find(".ProfileNav-item--tweets.is-active > a > span.ProfileNav-value").First()),
		URL:            url,
		Username:       doc.Find(".u-linkComplex-target").First().Text(),
		Website:        strings.TrimSpace(doc.Find(".ProfileHeaderCard-urlText.u-dir > a").First().AttrOr("title", "")),
	}, nil
}

func parseCount(sel *goquery.Selection) (i int) {
	if str, exists := sel.Attr("data-count"); exists {
		i, _ = strconv.Atoi(str)
	}
	return
}
