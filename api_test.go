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

func TestGetUserIDByScreenName(t *testing.T) {
	scraper := New()
	userID, err := scraper.GetUserIDByScreenName("Twitter")
	if err != nil {
		t.Errorf("getUserByScreenName() error = %v", err)
	}
	if userID == "" {
		t.Error("Expected non-empty user ID")
	}
}
