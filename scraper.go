package twitterscraper

import (
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Scraper object
type Scraper struct {
	client         *http.Client
	delay          int64
	guestToken     string
	guestCreatedAt time.Time
	includeReplies bool
	searchMode     SearchMode
	wg             sync.WaitGroup

	Cookie     string
	XCsrfToken string
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

// WithDelay add delay between API requests (in seconds)
func (s *Scraper) WithDelay(seconds int64) *Scraper {
	s.delay = seconds
	return s
}

// WithDelay wrapper for default Scraper
func WithDelay(seconds int64) *Scraper {
	return defaultScraper.WithDelay(seconds)
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

// cookie
func (s *Scraper) WithCookie(cookie string) *Scraper {
	s.Cookie = cookie
	return s
}

// x csrf token
func (s *Scraper) WithXCsrfToken(xcsrfToken string) *Scraper {
	s.XCsrfToken = xcsrfToken
	return s
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

// SetProxy set socks5 proxy in the format `HOST:PORT`
func (s *Scraper) SetSocks5Proxy(socks5 string) error {
	log.Println(socks5)
	if dialer, err := proxy.SOCKS5("tcp", socks5, nil, proxy.Direct); err != nil {
		return errors.New(fmt.Sprintf("can't connect to the socks5 proxy: %s, err: %s", socks5, err.Error()))
	} else {
		s.client = &http.Client{
			Transport: &http.Transport{
				Dial: dialer.Dial,
				// TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
				DialContext: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).DialContext,
			},
		}

		log.Println(s.client)
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
