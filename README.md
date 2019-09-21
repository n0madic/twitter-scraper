# Twitter Scraper

Golang implementation of python library <https://github.com/kennethreitz/twitter-scraper>

Twitter's API is annoying to work with, and has lots of limitations —
luckily their frontend (JavaScript) has it's own API, which I reverse-engineered.
No API rate limits. No tokens needed. No restrictions. Extremely fast.

You can use this library to get the text of any user's Tweets trivially.

## Usage

### Get tweets

```golang
package main

import (
    "fmt"
    twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
    for tweet := range twitterscraper.GetTweets("kennethreitz", 25) {
        if tweet.Error != nil {
            panic(tweet.Error)
        }
        fmt.Println(tweet.HTML)
    }
}
```

It appears you can ask for up to 25 pages of tweets reliably (~486 tweets).

### Get profile

```golang
package main

import (
    "fmt"
    twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
    profile, err := twitterscraper.GetProfile("kennethreitz")
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", profile)
}
```

## Installation

```shell
go get -u github.com/n0madic/twitter-scraper
```