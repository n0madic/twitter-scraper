package twitterscraper

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGetProfile(t *testing.T) {
	loc := time.FixedZone("UTC", 0)
	joined := time.Date(2007, 02, 20, 14, 35, 54, 0, loc)
	sample := Profile{
		Avatar:    "https://pbs.twimg.com/profile_images/1308010958862905345/-SGZioPb_normal.jpg",
		Banner:    "https://pbs.twimg.com/profile_banners/783214/1609475315",
		Biography: "What's happening?!",
		//	Birthday:   "March 21",
		IsPrivate:      false,
		IsVerified:     true,
		Joined:         &joined,
		Location:       "everywhere",
		Name:           "Twitter",
		PinnedTweetIDs: []string{},
		URL:            "https://twitter.com/Twitter",
		UserID:         "783214",
		Username:       "Twitter",
		Website:        "https://about.twitter.com/",
	}

	profile, err := GetProfile("Twitter")
	if err != nil {
		t.Error(err)
	}

	cmpOptions := cmp.Options{
		cmpopts.IgnoreFields(Profile{}, "FollowersCount"),
		cmpopts.IgnoreFields(Profile{}, "FollowingCount"),
		cmpopts.IgnoreFields(Profile{}, "FriendsCount"),
		cmpopts.IgnoreFields(Profile{}, "LikesCount"),
		cmpopts.IgnoreFields(Profile{}, "ListedCount"),
		cmpopts.IgnoreFields(Profile{}, "TweetsCount"),
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
	joined := time.Date(2009, 8, 12, 6, 18, 29, 0, loc)
	sample := Profile{
		Avatar:    "https://pbs.twimg.com/profile_images/1352282054256324610/_v3nslbW_normal.jpg",
		Banner:    "https://pbs.twimg.com/profile_banners/64958707/1551520603",
		Biography: "",
		//	Birthday:   "March 21",
		IsPrivate:      true,
		IsVerified:     false,
		Joined:         &joined,
		Location:       "",
		Name:           "saidah.jpg",
		PinnedTweetIDs: []string{},
		URL:            "https://twitter.com/sdhftrh",
		UserID:         "64958707",
		Username:       "sdhftrh",
		Website:        "https://youtu.be/0liuo2Q4bGo",
	}

	// some random private profile (found via google)
	profile, err := GetProfile("sdhftrh")
	if err != nil {
		t.Error(err)
	}

	cmpOptions := cmp.Options{
		cmpopts.IgnoreFields(Profile{}, "FollowersCount"),
		cmpopts.IgnoreFields(Profile{}, "FollowingCount"),
		cmpopts.IgnoreFields(Profile{}, "FriendsCount"),
		cmpopts.IgnoreFields(Profile{}, "LikesCount"),
		cmpopts.IgnoreFields(Profile{}, "ListedCount"),
		cmpopts.IgnoreFields(Profile{}, "TweetsCount"),
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

func TestGetProfileErrorSuspended(t *testing.T) {
	_, err := GetProfile("123")
	if err == nil {
		t.Error("Expected Error, got success")
	} else {
		if err.Error() != "Authorization: User has been suspended. (63)" {
			t.Errorf("Expected error 'Authorization: User has been suspended. (63)', got '%s'", err)
		}
	}
}

func TestGetProfileErrorNotFound(t *testing.T) {
	_, err := GetProfile("sample3123131")
	if err == nil {
		t.Error("Expected Error, got success")
	} else {
		if err.Error() != "Not found" {
			t.Errorf("Expected error 'Not found', got '%s'", err)
		}
	}
}
