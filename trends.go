package twitterscraper

import (
	"github.com/PuerkitoBio/goquery"
)

const trendsURL = "https://twitter.com/i/trends"

// GetTrends return list of trends.
func GetTrends() ([]string, error) {
	req, err := newRequest(trendsURL)
	if err != nil {
		return nil, err
	}

	htm, err := getHTMLFromJSON(req, "module_html")
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(htm)
	if err != nil {
		return nil, err
	}

	var trends []string
	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		if trend, ok := s.Attr("data-trend-name"); ok {
			trends = append(trends, trend)
		}
	})
	return trends, nil
}
