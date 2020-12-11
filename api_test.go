package twitterscraper

import (
	"testing"
)

func TestGetGuestToken(t *testing.T) {
	if err := GetGuestToken(); err != nil {
		t.Errorf("getGuestToken() error = %v", err)
	}
	if guestToken == "" {
		t.Error("Expected non-empty guestToken")
	}
}

func TestGetUserIDByScreenName(t *testing.T) {
	userID, err := GetUserIDByScreenName("Twitter")
	if err != nil {
		t.Errorf("getUserByScreenName() error = %v", err)
	}
	if userID == "" {
		t.Error("Expected non-empty user ID")
	}
}
