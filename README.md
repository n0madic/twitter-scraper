# Twitter Scraper

[![Go Reference](https://pkg.go.dev/badge/github.com/n0madic/twitter-scraper.svg)](https://pkg.go.dev/github.com/n0madic/twitter-scraper)

Twitter's API is annoying to work with, and has lots of limitations —
luckily their frontend (JavaScript) has it's own API, which I reverse-engineered.
No API rate limits. No tokens needed. No restrictions. Extremely fast.

You can use this library to get the text of any user's Tweets trivially.
## Prerequisites before downloading the Twitter Scraper
Install Go from the website

### Enable dependency tracking for your code
For enabling dependency tracking of the code you have to create a go.mod file
Run the command 
```
go mod init [path-to-module]
```
path-to-module needs to be the location of a remote repository from where Go can download the source code.
Typically, it's of the format
```
<prefix>/<descriptive-text>
```
Example: github.com/<project-name>

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
    scraper := twitterscraper.New()

    for tweet := range scraper.GetTweets(context.Background(), "Twitter", 50) {
        if tweet.Error != nil {
            panic(tweet.Error)
        }
        fmt.Println(tweet.Text)
    }
}
```

It appears you can ask for up to 50 tweets (limit ~3200 tweets).

### Get single tweet

```golang
package main

import (
    "fmt"

    twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
    scraper := twitterscraper.New()
    tweet, err := scraper.GetTweet("1328684389388185600")
    if err != nil {
        panic(err)
    }
    fmt.Println(tweet.Text)
}
```

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
    scraper := twitterscraper.New()
    for tweet := range scraper.SearchTweets(context.Background(),
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


#### Set search mode

```golang
scraper.SetSearchMode(twitterscraper.SearchLatest)
```

Options:

* `twitterscraper.SearchTop` - default mode
* `twitterscraper.SearchLatest` - live mode
* `twitterscraper.SearchPhotos` - image mode
* `twitterscraper.SearchVideos` - video mode
* `twitterscraper.SearchUsers` - user mode

### Get profile

```golang
package main

import (
    "fmt"
    twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
    scraper := twitterscraper.New()
    profile, err := scraper.GetProfile("Twitter")
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", profile)
}
```

### Search profiles by query

```golang
package main

import (
    "context"
    "fmt"
    twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
    scraper := twitterscraper.New().SetSearchMode(twitterscraper.SearchUsers)
    for profile := range scraper.SearchProfiles(context.Background(), "Twitter", 50) {
        if profile.Error != nil {
            panic(profile.Error)
        }
        fmt.Println(profile.Name)
    }
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
    scraper := twitterscraper.New()
    trends, err := scraper.GetTrends()
    if err != nil {
        panic(err)
    }
    fmt.Println(trends)
}
```

### Use cookie authentication

Some specified user tweets are protected that you must login and follow.
Cookie and xCsrfToken is optional.

```golang
scraper.WithCookie("twitter cookie after login")
scraper.WithXCsrfToken("twitter X-Csrf-Token after login")
```

### Use Proxy

Support HTTP(s) and SOCKS5 proxy

#### with HTTP

```golang
err := scraper.SetProxy("http://localhost:3128")
if err != nil {
    panic(err)
}
```

#### with SOCKS5

```golang
err := scraper.SetProxy("socks5://localhost:1080")
if err != nil {
    panic(err)
}
```

### Delay requests

Add delay between API requests (in seconds)

```golang
scraper.WithDelay(5)
```

### Load timeline with tweet replies

```golang
scraper.WithReplies(true)
```
