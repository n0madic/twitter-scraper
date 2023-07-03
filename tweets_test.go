package twitterscraper_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

var cmpOptions = cmp.Options{
	cmpopts.IgnoreFields(twitterscraper.Tweet{}, "Likes"),
	cmpopts.IgnoreFields(twitterscraper.Tweet{}, "Replies"),
	cmpopts.IgnoreFields(twitterscraper.Tweet{}, "Retweets"),
}

func TestGetTweets(t *testing.T) {
	count := 0
	maxTweetsNbr := 300
	dupcheck := make(map[string]bool)
	scraper := twitterscraper.New()
	err := scraper.LoginOpenAccount()
	if err != nil {
		t.Fatalf("LoginOpenAccount() error = %v", err)
	}
	for tweet := range scraper.GetTweets(context.Background(), "Twitter", maxTweetsNbr) {
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

func assertGetTweet(t *testing.T, expectedTweet *twitterscraper.Tweet) {
	actualTweet, err := testScraper.GetTweet(expectedTweet.ID)
	if err != nil {
		t.Error(err)
	} else if diff := cmp.Diff(expectedTweet, actualTweet, cmpOptions...); diff != "" {
		t.Error("Resulting tweet does not match the sample", diff)
	}
}

func TestGetTweetWithVideo(t *testing.T) {
	expectedTweet := twitterscraper.Tweet{
		ConversationID: "1328684389388185600",
		HTML:           "That thing you didn’t Tweet but wanted to but didn’t but got so close but then were like nah. <br><br>We have a place for that now—Fleets! <br><br>Rolling out to everyone starting today. <br><a href=\"https://t.co/auQAHXZMfH\"><img src=\"https://pbs.twimg.com/amplify_video_thumb/1328684333599756289/img/cP5KwbIXbGunNSBy.jpg\"/></a>",
		ID:             "1328684389388185600",
		Name:           "Twitter",
		PermanentURL:   "https://twitter.com/Twitter/status/1328684389388185600",
		Photos:         nil,
		Text:           "That thing you didn’t Tweet but wanted to but didn’t but got so close but then were like nah. \n\nWe have a place for that now—Fleets! \n\nRolling out to everyone starting today. https://t.co/auQAHXZMfH",
		TimeParsed:     time.Date(2020, 11, 17, 13, 0, 18, 0, time.FixedZone("UTC", 0)),
		Timestamp:      1605618018,
		UserID:         "783214",
		Username:       "Twitter",
		Videos: []twitterscraper.Video{{
			ID:      "1328684333599756289",
			Preview: "https://pbs.twimg.com/amplify_video_thumb/1328684333599756289/img/cP5KwbIXbGunNSBy.jpg",
			URL:     "https://video.twimg.com/amplify_video/1328684333599756289/vid/960x720/PcL8yv8KhgQ48Qpt.mp4?tag=13",
		}},
	}
	assertGetTweet(t, &expectedTweet)
}

func TestGetTweetWithMultiplePhotos(t *testing.T) {
	expectedTweet := twitterscraper.Tweet{
		ConversationID: "1390026628957417473",
		HTML:           `no bird too tall, no crop too short<br><br>introducing bigger and better images on iOS and Android, now available to everyone <br><a href="https://t.co/2buHfhfRAx"><img src="https://pbs.twimg.com/media/E0pd2L2XEAQ_gnn.jpg"/></a><br><img src="https://pbs.twimg.com/media/E0pd2hPXoAY9-TZ.jpg"/>`,
		ID:             "1390026628957417473",
		Name:           "Twitter",
		PermanentURL:   "https://twitter.com/Twitter/status/1390026628957417473",
		Photos: []twitterscraper.Photo{
			{ID: "1390026620472332292", URL: "https://pbs.twimg.com/media/E0pd2L2XEAQ_gnn.jpg"},
			{ID: "1390026626214371334", URL: "https://pbs.twimg.com/media/E0pd2hPXoAY9-TZ.jpg"},
		},
		Text:       "no bird too tall, no crop too short\n\nintroducing bigger and better images on iOS and Android, now available to everyone https://t.co/2buHfhfRAx",
		TimeParsed: time.Date(2021, 5, 5, 19, 32, 28, 0, time.FixedZone("UTC", 0)),
		Timestamp:  1620243148,
		UserID:     "783214",
		Username:   "Twitter",
	}
	assertGetTweet(t, &expectedTweet)
}

func TestGetTweetWithGIF(t *testing.T) {
	if os.Getenv("SKIP_AUTH_TEST") != "" {
		t.Skip("Skipping test due to environment variable")
	}
	expectedTweet := twitterscraper.Tweet{
		ConversationID: "1288540609310056450",
		GIFs: []twitterscraper.GIF{
			{
				ID:      "1288540582768517123",
				Preview: "https://pbs.twimg.com/tweet_video_thumb/EeHQ1UKXoAMVxWB.jpg",
				URL:     "https://video.twimg.com/tweet_video/EeHQ1UKXoAMVxWB.mp4",
			},
		},
		Hashtags:     []string{"CountdownToMars"},
		HTML:         `Like for liftoff! <a href="https://twitter.com/hashtag/CountdownToMars">#CountdownToMars</a> <br><a href="https://t.co/yLe331pHfY"><img src="https://pbs.twimg.com/tweet_video_thumb/EeHQ1UKXoAMVxWB.jpg"/></a>`,
		ID:           "1288540609310056450",
		Name:         "Twitter",
		PermanentURL: "https://twitter.com/Twitter/status/1288540609310056450",
		Text:         "Like for liftoff! #CountdownToMars https://t.co/yLe331pHfY",
		TimeParsed:   time.Date(2020, 7, 29, 18, 23, 15, 0, time.FixedZone("UTC", 0)),
		Timestamp:    1596046995,
		UserID:       "783214",
		Username:     "Twitter",
	}
	assertGetTweet(t, &expectedTweet)
}

func TestGetTweetWithPhotoAndGIF(t *testing.T) {
	if os.Getenv("SKIP_AUTH_TEST") != "" {
		t.Skip("Skipping test due to environment variable")
	}
	expectedTweet := twitterscraper.Tweet{
		ConversationID: "1580661436132757506",
		GIFs: []twitterscraper.GIF{
			{
				ID:      "1580661428335382531",
				Preview: "https://pbs.twimg.com/tweet_video_thumb/Fe-jMcIXkAMXK_W.jpg",
				URL:     "https://video.twimg.com/tweet_video/Fe-jMcIXkAMXK_W.mp4",
			},
		},
		HTML:         `a hit Tweet <br><a href="https://t.co/2C7cah4KzW"><img src="https://pbs.twimg.com/media/Fe-jMcGWQAAFWoG.jpg"/></a><br><img src="https://pbs.twimg.com/tweet_video_thumb/Fe-jMcIXkAMXK_W.jpg"/>`,
		ID:           "1580661436132757506",
		Name:         "Twitter",
		PermanentURL: "https://twitter.com/Twitter/status/1580661436132757506",
		Photos:       []twitterscraper.Photo{{ID: "1580661428326907904", URL: "https://pbs.twimg.com/media/Fe-jMcGWQAAFWoG.jpg"}},
		Text:         "a hit Tweet https://t.co/2C7cah4KzW",
		TimeParsed:   time.Date(2022, 10, 13, 20, 47, 8, 0, time.FixedZone("UTC", 0)),
		Timestamp:    1665694028,
		UserID:       "783214",
		Username:     "Twitter",
	}
	assertGetTweet(t, &expectedTweet)
}

func TestTweetMentions(t *testing.T) {
	sample := []twitterscraper.Mention{{
		ID:       "7018222",
		Username: "davidmcraney",
		Name:     "David McRaney",
	}}
	tweet, err := testScraper.GetTweet("1554522888904101890")
	if err != nil {
		t.Error(err)
	} else {
		if diff := cmp.Diff(sample, tweet.Mentions, cmpOptions...); diff != "" {
			t.Error("Resulting tweet does not match the sample", diff)
		}
	}
}

func TestQuotedAndReply(t *testing.T) {
	sample := &twitterscraper.Tweet{
		ConversationID: "1237110546383724547",
		HTML:           "The Easiest Problem Everyone Gets Wrong <br><br>[new video] --&gt; <a href=\"https://youtu.be/ytfCdqWhmdg\">https://t.co/YdaeDYmPAU</a> <br><a href=\"https://t.co/iKu4Xs6o2V\"><img src=\"https://pbs.twimg.com/media/ESsZa9AXgAIAYnF.jpg\"/></a>",
		ID:             "1237110546383724547",
		Likes:          485,
		Name:           "Vsauce2",
		PermanentURL:   "https://twitter.com/VsauceTwo/status/1237110546383724547",
		Photos: []twitterscraper.Photo{{
			ID:  "1237110473486729218",
			URL: "https://pbs.twimg.com/media/ESsZa9AXgAIAYnF.jpg",
		}},
		Replies:    12,
		Retweets:   18,
		Text:       "The Easiest Problem Everyone Gets Wrong \n\n[new video] --&gt; https://t.co/YdaeDYmPAU https://t.co/iKu4Xs6o2V",
		TimeParsed: time.Date(2020, 0o3, 9, 20, 18, 33, 0, time.FixedZone("UTC", 0)),
		Timestamp:  1583785113,
		URLs:       []string{"https://youtu.be/ytfCdqWhmdg"},
		UserID:     "978944851",
		Username:   "VsauceTwo",
	}
	tweet, err := testScraper.GetTweet("1237110897597976576")
	if err != nil {
		t.Error(err)
	} else {
		if !tweet.IsQuoted {
			t.Error("IsQuoted must be True")
		}
		if diff := cmp.Diff(sample, tweet.QuotedStatus, cmpOptions...); diff != "" {
			t.Error("Resulting quote does not match the sample", diff)
		}
	}
	tweet, err = testScraper.GetTweet("1237111868445134850")
	if err != nil {
		t.Error(err)
	} else {
		if !tweet.IsReply {
			t.Error("IsReply must be True")
		}
		if diff := cmp.Diff(sample, tweet.InReplyToStatus, cmpOptions...); diff != "" {
			t.Error("Resulting reply does not match the sample", diff)
		}
	}

}
func TestRetweet(t *testing.T) {
	sample := &twitterscraper.Tweet{
		ConversationID: "1359151057872580612",
		HTML:           "We’ve seen an increase in attacks against Asian communities and individuals around the world. It’s important to know that this isn’t new; throughout history, Asians have experienced violence and exclusion. However, their diverse lived experiences have largely been overlooked.",
		ID:             "1359151057872580612",
		IsSelfThread:   false,
		Likes:          6683,
		Name:           "Twitter Together",
		PermanentURL:   "https://twitter.com/TwitterTogether/status/1359151057872580612",
		Replies:        456,
		Retweets:       1495,
		Text:           "We’ve seen an increase in attacks against Asian communities and individuals around the world. It’s important to know that this isn’t new; throughout history, Asians have experienced violence and exclusion. However, their diverse lived experiences have largely been overlooked.",
		TimeParsed:     time.Date(2021, 02, 9, 14, 43, 58, 0, time.FixedZone("UTC", 0)),
		Timestamp:      1612881838,
		UserID:         "773578328498372608",
		Username:       "TwitterTogether",
	}
	tweet, err := testScraper.GetTweet("1362849141248974853")
	if err != nil {
		t.Error(err)
	} else {
		if !tweet.IsRetweet {
			t.Error("IsRetweet must be True")
		}
		if diff := cmp.Diff(sample, tweet.RetweetedStatus, cmpOptions...); diff != "" {
			t.Error("Resulting retweet does not match the sample", diff)
		}
	}
}

func TestTweetViews(t *testing.T) {
	sample := &twitterscraper.Tweet{
		HTML:         "Replies and likes don’t tell the whole story. We’re making it easier to tell *just* how many people have seen your Tweets with the addition of view counts, shown right next to likes. Now on iOS and Android, web coming soon.",
		ID:           "1606055187348688896",
		Likes:        2839,
		Name:         "Twitter Support",
		PermanentURL: "https://twitter.com/TwitterSupport/status/1606055187348688896",
		Replies:      3427,
		Retweets:     783,
		Text:         "Replies and likes don’t tell the whole story. We’re making it easier to tell *just* how many people have seen your Tweets with the addition of view counts, shown right next to likes. Now on iOS and Android, web coming soon.",
		TimeParsed:   time.Date(2022, 12, 22, 22, 32, 50, 0, time.FixedZone("UTC", 0)),
		Timestamp:    1612881838,
		UserID:       "17874544",
		Username:     "TwitterSupport",
		Views:        3189278,
	}
	tweet, err := testScraper.GetTweet("1606055187348688896")
	if err != nil {
		t.Error(err)
	} else {
		if tweet.Views < sample.Views {
			t.Error("Views must be greater than or equal to the sample")
		}
	}
}

func TestTweetThread(t *testing.T) {
	if os.Getenv("SKIP_AUTH_TEST") != "" {
		t.Skip("Skipping test due to environment variable")
	}
	tweet, err := testScraper.GetTweet("1665602315745673217")
	if err != nil {
		t.Fatal(err)
	} else {
		if !tweet.IsSelfThread {
			t.Error("IsSelfThread must be True")
		}
		if len(tweet.Thread) != 7 {
			t.Error("Thread length must be 7")
		}
	}
}
