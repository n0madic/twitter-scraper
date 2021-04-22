package twitterscraper

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGetTweets(t *testing.T) {
	count := 0
	maxTweetsNbr := 300
	dupcheck := make(map[string]bool)
	for tweet := range GetTweets(context.Background(), "Twitter", maxTweetsNbr) {
		if tweet.Error != nil {
			t.Error(tweet.Error)
		} else {
			count++
			if tweet.ID == "" {
				t.Error("Expected tweet ID is empty")
			} else {
				if dupcheck[tweet.ID] {
					t.Errorf("Detect duplicated tweet ID: %s", tweet.ID)
				} else {
					dupcheck[tweet.ID] = true
				}
			}
			if tweet.UserID == "" {
				t.Error("Expected tweet UserID is empty")
			}
			if tweet.Username == "" {
				t.Error("Expected tweet Username is empty")
			}
			if tweet.PermanentURL == "" {
				t.Error("Expected tweet PermanentURL is empty")
			}
			if tweet.Text == "" {
				t.Error("Expected tweet Text is empty")
			}
			if tweet.TimeParsed.IsZero() {
				t.Error("Expected tweet TimeParsed is zero")
			}
			if tweet.Timestamp == 0 {
				t.Error("Expected tweet Timestamp is greater than zero")
			}
			for _, video := range tweet.Videos {
				if video.ID == "" {
					t.Error("Expected tweet video ID is empty")
				}
				if video.Preview == "" {
					t.Error("Expected tweet video Preview is empty")
				}
				if video.URL == "" {
					t.Error("Expected tweet video URL is empty")
				}
			}
		}
	}
	if count != maxTweetsNbr {
		t.Errorf("Expected tweets count=%v, got: %v", maxTweetsNbr, count)
	}
}

func TestGetTweet(t *testing.T) {
	sample := Tweet{
		HTML:         "That thing you didn’t Tweet but wanted to but didn’t but got so close but then were like nah. <br><br>We have a place for that now—Fleets! <br><br>Rolling out to everyone starting today. <br><a href=\"https://t.co/auQAHXZMfH\"><img src=\"https://pbs.twimg.com/amplify_video_thumb/1328684333599756289/img/cP5KwbIXbGunNSBy.jpg\"/></a>",
		ID:           "1328684389388185600",
		PermanentURL: "https://twitter.com/Twitter/status/1328684389388185600",
		Photos:       []string{"https://pbs.twimg.com/amplify_video_thumb/1328684333599756289/img/cP5KwbIXbGunNSBy.jpg"},
		Text:         "That thing you didn’t Tweet but wanted to but didn’t but got so close but then were like nah. \n\nWe have a place for that now—Fleets! \n\nRolling out to everyone starting today. https://t.co/auQAHXZMfH",
		TimeParsed:   time.Date(2020, 11, 17, 13, 0, 18, 0, time.FixedZone("UTC", 0)),
		Timestamp:    1605618018,
		UserID:       "783214",
		Username:     "Twitter",
		Videos: []Video{{
			ID:      "1328684333599756289",
			Preview: "https://pbs.twimg.com/amplify_video_thumb/1328684333599756289/img/cP5KwbIXbGunNSBy.jpg",
			URL:     "https://video.twimg.com/amplify_video/1328684333599756289/vid/480x360/Qh70ELAcq-N2RYmZ.mp4?tag=13",
		}},
	}
	tweet, err := defaultScraper.GetTweet("1328684389388185600")
	if err != nil {
		t.Error(err)
	} else {
		cmpOptions := cmp.Options{
			cmpopts.IgnoreFields(Tweet{}, "Likes"),
			cmpopts.IgnoreFields(Tweet{}, "Replies"),
			cmpopts.IgnoreFields(Tweet{}, "Retweets"),
		}
		if diff := cmp.Diff(sample, *tweet, cmpOptions...); diff != "" {
			t.Error("Resulting tweet does not match the sample", diff)
		}
	}
}

func TestRetweet(t *testing.T) {
	sample := Retweet{
		ID:         "1359151057872580612",
		TimeParsed: time.Date(2021, 02, 9, 14, 43, 58, 0, time.FixedZone("UTC", 0)),
		Timestamp:  1612881838,
		UserID:     "773578328498372608",
		Username:   "TwitterTogether",
	}
	tweet, err := defaultScraper.GetTweet("1362849141248974853")
	if err != nil {
		t.Error(err)
	} else {
		if diff := cmp.Diff(sample, tweet.Retweet); diff != "" {
			t.Error("Resulting retweet does not match the sample", diff)
		}
	}
}
