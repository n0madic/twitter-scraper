package twitterscraper

import (
	"crypto/tls"
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
	guestCreatedAt time.Time
	includeReplies bool
	searchMode     SearchMode
}

// SearchMode type
type SearchMode int

const (
	// SearchTop - default mode
	SearchTop SearchMode = iota
	// SearchLatest - live mode
	SearchLatest
	// SearchPhotos - image mode
	SearchPhotos
	// SearchVideos - video mode
	SearchVideos
	// SearchUsers - user mode
	SearchUsers
)

var defaultScraper *Scraper

// New creates a Scraper object
func New() *Scraper {
	return &Scraper{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// SetSearchMode switcher
func (s *Scraper) SetSearchMode(mode SearchMode) *Scraper {
	s.searchMode = mode
	return s
}

// SetSearchMode wrapper for default Scraper
func SetSearchMode(mode SearchMode) *Scraper {
	return defaultScraper.SetSearchMode(mode)
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
	if !strings.HasPrefix(proxy, "http") {
		return errors.New("only support http(s) protocol")
	}
	urlproxy, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	s.client = &http.Client{
		Transport: &http.Transport{
			Proxy:        http.ProxyURL(urlproxy),
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
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
