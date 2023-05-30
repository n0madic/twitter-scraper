package twitterscraper_test

import (
	"testing"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

func TestGetTrends(t *testing.T) {
	scraper := twitterscraper.New()
	trends, err := scraper.GetTrends()
	if err != nil {
		t.Fatal(err)
	}

	if len(trends) != 20 {
		t.Errorf("Expected 20 trends, got %d: %#v", len(trends), trends)
	}

	for _, trend := range trends {
		if trend == "" {
			t.Error("Expected trend is empty")
		}
	}
}
