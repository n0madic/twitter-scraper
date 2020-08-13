package twitterscraper

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGetProfile(t *testing.T) {
	// loc := time.FixedZone("UTC", 0)
	// joined := time.Date(2007, 02, 20, 6, 35, 0, 0, loc)
	sample := Profile{
		Avatar: "https://pbs.twimg.com/profile_images/1293262551057420293/4QMAmava_normal.jpg",
		// Banner:     "https://pbs.twimg.com/profile_banners/783214/1596041768/1500x500",
		Biography: "#BlackLivesMatter\n#BlackTransLivesMatter",
		// Birthday:   "March 21",
		IsPrivate:  false,
		IsVerified: true,
		// Joined:     &joined,
		Location: "everywhere",
		Name:     "Twitter",
		URL:      "https://twitter.com/Twitter",
		// UserID:     "783214",
		Username: "Twitter",
		Website:  "about.twitter.com",
	}

	profile, err := GetProfile("Twitter")
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
	// if profile.LikesCount == 0 {
	// 	t.Error("Expected LikesCount is greater than zero")
	// }
	if profile.TweetsCount == 0 {
		t.Error("Expected TweetsCount is greater than zero")
	}
}
