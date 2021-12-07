package twitterscraper_test

import (
	"testing"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

func TestGetGuestToken(t *testing.T) {
	scraper := twitterscraper.New()
	if err := scraper.GetGuestToken(); err != nil {
		t.Errorf("getGuestToken() error = %v", err)
	}
	if !scraper.IsGuestToken() {
		t.Error("Expected non-empty guestToken")
	}
}
