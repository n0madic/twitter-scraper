package twitterscraper_test

import (
	"testing"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

func TestGetTrends(t *testing.T) {
	trends, err := twitterscraper.GetTrends()
	if err != nil {
		t.Error(err)
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
