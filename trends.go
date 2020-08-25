package twitterscraper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const trendsURL = "https://mobile.twitter.com/trends"

// GetTrends return list of trends.
func GetTrends() ([]string, error) {
	req, err := http.NewRequest("GET", trendsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Language", "en-US")

	resp, err := http.DefaultClient.Do(req)
	if resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var trends []string
	doc.Find("li.topic").Each(func(i int, s *goquery.Selection) {
		trends = append(trends, strings.TrimSpace(s.Text()))
	})
	return trends, nil
}
