package twitterscraper

import (
	"strconv"
)

type result struct {
	Typename string `json:"__typename"`
	Core     struct {
		UserResults struct {
			Result struct {
				IsBlueVerified bool       `json:"is_blue_verified"`
				Legacy         legacyUser `json:"legacy"`
			} `json:"result"`
		} `json:"user_results"`
	} `json:"core"`
	Views struct {
		Count string `json:"count"`
	} `json:"views"`
	NoteTweet struct {
		NoteTweetResults struct {
			Result struct {
				Text string `json:"text"`
			} `json:"result"`
		} `json:"note_tweet_results"`
	} `json:"note_tweet"`
	QuotedStatusResult struct {
		Result *result `json:"result"`
	} `json:"quoted_status_result"`
	Legacy legacyTweet `json:"legacy"`
}

func (result *result) parse() *Tweet {
	if result.NoteTweet.NoteTweetResults.Result.Text != "" {
		result.Legacy.FullText = result.NoteTweet.NoteTweetResults.Result.Text
	}
	tw := parseLegacyTweet(&result.Core.UserResults.Result.Legacy, &result.Legacy)
	if tw == nil {
		return nil
	}
	if tw.Views == 0 && result.Views.Count != "" {
		tw.Views, _ = strconv.Atoi(result.Views.Count)
	}
	if result.QuotedStatusResult.Result != nil {
		tw.QuotedStatus = result.QuotedStatusResult.Result.parse()
	}
	return tw
}

type entry struct {
	Content struct {
		CursorType string `json:"cursorType"`
		Value      string `json:"value"`
		Items      []struct {
			Item struct {
				ItemContent struct {
					TweetDisplayType string `json:"tweetDisplayType"`
					TweetResults     struct {
						Result result `json:"result"`
					} `json:"tweet_results"`
				} `json:"itemContent"`
			} `json:"item"`
		} `json:"items"`
		ItemContent struct {
			TweetDisplayType string `json:"tweetDisplayType"`
			TweetResults     struct {
				Result result `json:"result"`
			} `json:"tweet_results"`
			UserDisplayType string `json:"userDisplayType"`
			UserResults     struct {
				Result struct {
					RestID string     `json:"rest_id"`
					Legacy legacyUser `json:"legacy"`
				} `json:"result"`
			} `json:"user_results"`
		} `json:"itemContent"`
	} `json:"content"`
}

// timeline v2 JSON object
type timelineV2 struct {
	Data struct {
		User struct {
			Result struct {
				TimelineV2 struct {
					Timeline struct {
						Instructions []struct {
							Entries []entry `json:"entries"`
							Entry   entry   `json:"entry"`
							Type    string  `json:"type"`
						} `json:"instructions"`
					} `json:"timeline"`
				} `json:"timeline_v2"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

func (timeline *timelineV2) parseTweets() ([]*Tweet, string) {
	var cursor string
	var tweets []*Tweet
	for _, instruction := range timeline.Data.User.Result.TimelineV2.Timeline.Instructions {
		for _, entry := range instruction.Entries {
			if entry.Content.CursorType == "Bottom" {
				cursor = entry.Content.Value
				continue
			}
			if entry.Content.ItemContent.TweetResults.Result.Typename == "Tweet" {
				if tweet := entry.Content.ItemContent.TweetResults.Result.parse(); tweet != nil {
					tweets = append(tweets, tweet)
				}
			}
		}
	}
	return tweets, cursor
}

type threadedConversation struct {
	Data struct {
		ThreadedConversationWithInjectionsV2 struct {
			Instructions []struct {
				Type    string  `json:"type"`
				Entries []entry `json:"entries"`
				Entry   entry   `json:"entry"`
			} `json:"instructions"`
		} `json:"threaded_conversation_with_injections_v2"`
	} `json:"data"`
}

func (conversation *threadedConversation) parse() []*Tweet {
	var tweets []*Tweet
	for _, instruction := range conversation.Data.ThreadedConversationWithInjectionsV2.Instructions {
		for _, entry := range instruction.Entries {
			if entry.Content.ItemContent.TweetResults.Result.Typename == "Tweet" {
				if tweet := entry.Content.ItemContent.TweetResults.Result.parse(); tweet != nil {
					if entry.Content.ItemContent.TweetDisplayType == "SelfThread" {
						tweet.IsSelfThread = true
					}
					tweets = append(tweets, tweet)
				}
			}
			for _, item := range entry.Content.Items {
				if item.Item.ItemContent.TweetResults.Result.Typename == "Tweet" {
					if tweet := item.Item.ItemContent.TweetResults.Result.parse(); tweet != nil {
						if item.Item.ItemContent.TweetDisplayType == "SelfThread" {
							tweet.IsSelfThread = true
						}
						tweets = append(tweets, tweet)
					}
				}
			}
		}
	}
	for _, tweet := range tweets {
		if tweet.InReplyToStatusID != "" {
			for _, parentTweet := range tweets {
				if parentTweet.ID == tweet.InReplyToStatusID {
					tweet.InReplyToStatus = parentTweet
					break
				}
			}
		}
		if tweet.IsSelfThread && tweet.ConversationID == tweet.ID {
			for _, childTweet := range tweets {
				if childTweet.IsSelfThread && childTweet.ID != tweet.ID {
					tweet.Thread = append(tweet.Thread, childTweet)
				}
			}
			if len(tweet.Thread) == 0 {
				tweet.IsSelfThread = false
			}
		}
	}
	return tweets
}
