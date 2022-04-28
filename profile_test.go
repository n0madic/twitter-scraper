package twitterscraper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

func TestGetProfile(t *testing.T) {
	loc := time.FixedZone("UTC", 0)
	joined := time.Date(2010, 01, 18, 8, 49, 30, 0, loc)
	sample := twitterscraper.Profile{
		Avatar:    "https://pbs.twimg.com/profile_images/436075027193004032/XlDa2oaz_normal.jpeg",
		Banner:    "https://pbs.twimg.com/profile_banners/106037940/1541084318",
		Biography: "nothing",
		//	Birthday:   "March 21",
		IsPrivate:      false,
		IsVerified:     false,
		Joined:         &joined,
		Location:       "Ukraine",
		Name:           "Nomadic",
		PinnedTweetIDs: []string{},
		URL:            "https://twitter.com/nomadic_ua",
		UserID:         "106037940",
		Username:       "nomadic_ua",
		Website:        "https://nomadic.name",
	}

	profile, err := twitterscraper.GetProfile("nomadic_ua")
	if err != nil {
		t.Error(err)
	}

	cmpOptions := cmp.Options{
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "FollowersCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "FollowingCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "FriendsCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "LikesCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "ListedCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "TweetsCount"),
	}
	if diff := cmp.Diff(sample, profile, cmpOptions...); diff != "" {
		t.Error("Resulting profile does not match the sample", diff)
	}

	if profile.FollowersCount == 0 {
		t.Error("Expected FollowersCount is greater than zero")
	}
	if profile.FollowingCount == 0 {
		t.Error("Expected FollowingCount is greater than zero")
	}
	if profile.LikesCount == 0 {
		t.Error("Expected LikesCount is greater than zero")
	}
	if profile.TweetsCount == 0 {
		t.Error("Expected TweetsCount is greater than zero")
	}
}

func TestGetProfilePrivate(t *testing.T) {
	loc := time.FixedZone("UTC", 0)
	joined := time.Date(2020, 1, 26, 0, 3, 5, 0, loc)
	sample := twitterscraper.Profile{
		Avatar:    "https://pbs.twimg.com/profile_images/1222218816484020224/ik9P1QZt_normal.jpg",
		Banner:    "",
		Biography: "Quien es mas macho",
		//	Birthday:   "March 21",
		IsPrivate:      true,
		IsVerified:     false,
		Joined:         &joined,
		Location:       "",
		Name:           "private account",
		PinnedTweetIDs: []string{},
		URL:            "https://twitter.com/tomdumont",
		UserID:         "1221221876849995777",
		Username:       "tomdumont",
		Website:        "",
	}

	// some random private profile (found via google)
	profile, err := twitterscraper.GetProfile("tomdumont")
	if err != nil {
		t.Error(err)
	}

	cmpOptions := cmp.Options{
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "FollowersCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "FollowingCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "FriendsCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "LikesCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "ListedCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "TweetsCount"),
	}
	if diff := cmp.Diff(sample, profile, cmpOptions...); diff != "" {
		t.Error("Resulting profile does not match the sample", diff)
	}

	if profile.FollowingCount == 0 {
		t.Error("Expected FollowingCount is greater than zero")
	}
	if profile.LikesCount == 0 {
		t.Error("Expected LikesCount is greater than zero")
	}
	if profile.TweetsCount == 0 {
		t.Error("Expected TweetsCount is greater than zero")
	}
}

func TestGetProfileErrorSuspended(t *testing.T) {
	_, err := twitterscraper.GetProfile("123")
	if err == nil {
		t.Error("Expected Error, got success")
	} else {
		if err.Error() != "Authorization: User has been suspended. (63)" {
			t.Errorf("Expected error 'Authorization: User has been suspended. (63)', got '%s'", err)
		}
	}
}

func TestGetProfileErrorNotFound(t *testing.T) {
	neUser := "sample3123131"
	expectedError := fmt.Sprintf("User '%s' not found", neUser)
	_, err := twitterscraper.GetProfile(neUser)
	if err == nil {
		t.Error("Expected Error, got success")
	} else {
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, err)
		}
	}
}

func TestGetUserIDByScreenName(t *testing.T) {
	scraper := twitterscraper.New()
	userID, err := scraper.GetUserIDByScreenName("Twitter")
	if err != nil {
		t.Errorf("getUserByScreenName() error = %v", err)
	}
	if userID == "" {
		t.Error("Expected non-empty user ID")
	}
}
