package twitterscraper

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

// Scraper object
type Scraper struct {
	bearerToken    string
	client         *http.Client
	delay          int64
	guestToken     string
	guestCreatedAt time.Time
	includeReplies bool
	isLogged       bool
	searchMode     SearchMode
	wg             sync.WaitGroup
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

// default http client timeout
const DefaultClientTimeout = 10 * time.Second

// New creates a Scraper object
func New() *Scraper {
	jar, _ := cookiejar.New(nil)
	return &Scraper{
		bearerToken: bearerToken,
		client: &http.Client{
			Jar:     jar,
			Timeout: DefaultClientTimeout,
		},
	}
}

func (s *Scraper) setBearerToken(token string) {
	s.bearerToken = token
	s.guestToken = ""
}

// IsGuestToken check if guest token not empty
func (s *Scraper) IsGuestToken() bool {
	return s.guestToken != ""
}

// SetSearchMode switcher
func (s *Scraper) SetSearchMode(mode SearchMode) *Scraper {
	s.searchMode = mode
	return s
}

// WithDelay add delay between API requests (in seconds)
func (s *Scraper) WithDelay(seconds int64) *Scraper {
	s.delay = seconds
	return s
}

// WithReplies enable/disable load timeline with tweet replies
func (s *Scraper) WithReplies(b bool) *Scraper {
	s.includeReplies = b
	return s
}

// client timeout
func (s *Scraper) WithClientTimeout(timeout time.Duration) *Scraper {
	s.client.Timeout = timeout
	return s
}

// SetProxy
// set http proxy in the format `http://HOST:PORT`
// set socket proxy in the format `socks5://HOST:PORT`
func (s *Scraper) SetProxy(proxyAddr string) error {
	if strings.HasPrefix(proxyAddr, "http") {
		urlproxy, err := url.Parse(proxyAddr)
		if err != nil {
			return err
		}
		s.client.Transport = &http.Transport{
			Proxy:        http.ProxyURL(urlproxy),
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			DialContext: (&net.Dialer{
				Timeout: s.client.Timeout,
			}).DialContext,
		}
		return nil
	}
	if strings.HasPrefix(proxyAddr, "socks5") {
		baseDialer := &net.Dialer{
			Timeout:   s.client.Timeout,
			KeepAlive: s.client.Timeout,
		}
		socksHostPort := strings.ReplaceAll(proxyAddr, "socks5://", "")
		dialSocksProxy, err := proxy.SOCKS5("tcp", socksHostPort, nil, baseDialer)
		if err != nil {
			return errors.New("error creating socks5 proxy :" + err.Error())
		}
		if contextDialer, ok := dialSocksProxy.(proxy.ContextDialer); ok {
			dialContext := contextDialer.DialContext
			s.client.Transport = &http.Transport{
				DialContext: dialContext,
			}
		} else {
			return errors.New("failed type assertion to DialContext")
		}
		return nil
	}
	return errors.New("only support http(s) or socks5 protocol")
}
