package twitterscraper

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Scraper object
type Scraper struct {
	client         *http.Client
	guestToken     string
	includeReplies bool
}

var defaultScraper Scraper

// New creates a Scraper object
func New() Scraper {
	return Scraper{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// WithReplies enable/disable load timeline with tweet replies
func (s *Scraper) WithReplies(b bool) *Scraper {
	s.includeReplies = b
	return s
}

// WithReplies wrapper for default Scraper
func WithReplies(b bool) *Scraper {
	return defaultScraper.WithReplies(b)
}

// SetProxy set http proxy in the format `http://HOST:PORT`
func (s *Scraper) SetProxy(proxy string) error {
	if !strings.HasPrefix(proxy, "http://") {
		return errors.New("only support http protocol")
	}
	urlproxy, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	s.client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
		},
	}
	return nil
}

// SetProxy wrapper for default Scraper
func SetProxy(proxy string) error {
	return defaultScraper.SetProxy(proxy)
}

func init() {
	defaultScraper = New()
}
