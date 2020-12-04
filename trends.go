package twitterscraper

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const trendsURL = "https://mobile.twitter.com/trends"

// GetTrends return list of trends.
func GetTrends() ([]string, error) {
	client := http.DefaultClient
	if HTTPProxy != nil {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(HTTPProxy),
				DialContext: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).DialContext,
			},
		}
	}

	req, err := http.NewRequest("GET", trendsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Language", "en-US")

	resp, err := client.Do(req)
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
