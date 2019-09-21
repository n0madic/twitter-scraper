package twitterscraper

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGetProfile(t *testing.T) {
	joined := time.Unix(1245860880, 0)
	sample := Profile{
		Avatar:    "https://pbs.twimg.com/profile_images/1174251196758003713/IjXgRpzL_400x400.jpg",
		Biography: "Designer of interfaces & abstractions.",
		Birthday:  "1988",
		Joined:    &joined,
		Location:  "Eden, Earth, Milky Way",
		Name:      "☿ Kenneth Reitz",
		URL:       "https://twitter.com/kennethreitz",
		Username:  "kennethreitz",
		Website:   "https://kennethreitz.org/values",
	}

	profile, err := GetProfile("kennethreitz")
	if err != nil {
		t.Error(err)
	}

	var cmpOptions = cmp.Options{
		cmpopts.IgnoreFields(Profile{}, "FollowersCount"),
		cmpopts.IgnoreFields(Profile{}, "FollowingCount"),
		cmpopts.IgnoreFields(Profile{}, "LikesCount"),
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
