# Twitter Scraper

Twitter's API is annoying to work with, and has lots of limitations —
luckily their frontend (JavaScript) has it's own API, which I reverse-engineered.
No API rate limits. No tokens needed. No restrictions. Extremely fast.

You can use this library to get the text of any user's Tweets trivially.

## Installation

```shell
go get -u github.com/n0madic/twitter-scraper
```

## Usage

### Get user tweets

```golang
package main

import (
    "context"
    "fmt"
    twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
    for tweet := range twitterscraper.GetTweets(context.Background(), "Twitter", 50) {
        if tweet.Error != nil {
            panic(tweet.Error)
        }
        fmt.Println(tweet.Text)
    }
}
```

It appears you can ask for up to 50 tweets (limit ~3200 tweets).

### Search tweets by query standard operators

Tweets containing “twitter” and “scraper” and “data“, filtering out retweets:

```golang
package main

import (
    "context"
    "fmt"
    twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
    for tweet := range twitterscraper.SearchTweets(context.Background(),
        "twitter scraper data -filter:retweets", 50) {
        if tweet.Error != nil {
            panic(tweet.Error)
        }
        fmt.Println(tweet.Text)
    }
}
```

The search ends if we have 50 tweets.

See [Rules and filtering](https://developer.twitter.com/en/docs/tweets/rules-and-filtering/overview/standard-operators) for build standard queries.

### Get profile

```golang
package main

import (
    "fmt"
    twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
    profile, err := twitterscraper.GetProfile("Twitter")
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", profile)
}
```

### Get trends

```golang
package main

import (
    "fmt"
    twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
    trends, err := twitterscraper.GetTrends()
    if err != nil {
        panic(err)
    }
    fmt.Println(trends)
}
```

### Use http proxy

```golang
twitterscraper.SetProxy("http://localhost:3128")
```

### Load timeline with tweet replies

```golang
twitterscraper.IncludeReplies = true
```
