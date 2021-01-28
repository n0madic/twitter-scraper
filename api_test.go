package twitterscraper

import (
	"testing"
)

func TestGetGuestToken(t *testing.T) {
	scraper := New()
	if err := scraper.GetGuestToken(); err != nil {
		t.Errorf("getGuestToken() error = %v", err)
	}
	if scraper.guestToken == "" {
		t.Error("Expected non-empty guestToken")
	}
}
