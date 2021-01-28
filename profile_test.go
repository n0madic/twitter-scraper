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
		Avatar:    "https://pbs.twimg.com/profile_images/1354479643882004483/Btnfm47p_normal.jpg",
		Banner:    "https://pbs.twimg.com/profile_banners/783214/1611770459",
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
	joined := time.Date(2020, 1, 26, 0, 3, 5, 0, loc)
	sample := Profile{
		Avatar:    "https://pbs.twimg.com/profile_images/1222218816484020224/ik9P1QZt_normal.jpg",
		Banner:    "",
		Biography: "private account",
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
	profile, err := GetProfile("tomdumont")
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
