package twitterscraper

import (
	"context"
	"testing"
)

func TestFetchSearchCursor(t *testing.T) {
	scraper := New()
	maxTweetsNbr := 150
	tweetsNbr := 0
	nextCursor := ""
	for tweetsNbr < maxTweetsNbr {
		tweets, cursor, err := scraper.FetchSearchTweets("twitter", maxTweetsNbr, nextCursor)
		if err != nil {
			t.Fatal(err)
		}
		if cursor == "" {
			t.Fatal("Expected search cursor is not empty")
		}
		tweetsNbr += len(tweets)
		nextCursor = cursor
	}
}

func TestGetSearchTweets(t *testing.T) {
	count := 0
	maxTweetsNbr := 250
	dupcheck := make(map[string]bool)
	for tweet := range SearchTweets(context.Background(), "twitter -filter:retweets", maxTweetsNbr) {
		if tweet.Error != nil {
			t.Error(tweet.Error)
		} else {
			count++
			if tweet.ID == "" {
				t.Error("Expected tweet ID is not empty")
			} else {
				if dupcheck[tweet.ID] {
					t.Errorf("Detect duplicated tweet ID: %s", tweet.ID)
				} else {
					dupcheck[tweet.ID] = true
				}
			}
			if tweet.PermanentURL == "" {
				t.Error("Expected tweet PermanentURL is not empty")
			}
			if tweet.IsRetweet {
				t.Error("Expected tweet IsRetweet is false")
			}
			if tweet.Text == "" {
				t.Error("Expected tweet Text is not empty")
			}
		}
	}

	if count != maxTweetsNbr {
		t.Errorf("Expected tweets count=%v, got: %v", maxTweetsNbr, count)
	}
}
