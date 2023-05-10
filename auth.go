package twitterscraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

const (
	loginURL     = "https://api.twitter.com/1.1/onboarding/task.json"
	logoutURL    = "https://api.twitter.com/1.1/account/logout.json"
	bearerToken2 = "AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"
)

type (
	flow struct {
		Errors []struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
		FlowToken string `json:"flow_token"`
		Status    string `json:"status"`
		Subtasks  []struct {
			SubtaskID string `json:"subtask_id"`
		} `json:"subtasks"`
	}

	verifyCredentials struct {
		Errors []struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
	}
)

func (s *Scraper) getFlowToken(data map[string]interface{}) (string, error) {
	headers := http.Header{
		"Authorization":             []string{"Bearer " + s.bearerToken},
		"Content-Type":              []string{"application/json"},
		"User-Agent":                []string{"Mozilla/5.0 (Linux; Android 11; Nokia G20) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.88 Mobile Safari/537.36"},
		"X-Guest-Token":             []string{s.guestToken},
		"X-Twitter-Auth-Type":       []string{"OAuth2Client"},
		"X-Twitter-Active-User":     []string{"yes"},
		"X-Twitter-Client-Language": []string{"en"},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", loginURL, bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}
	req.Header = headers

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var info flow
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return "", err
	}

	if len(info.Errors) > 0 {
		return "", fmt.Errorf("auth error (%d): %v", info.Errors[0].Code, info.Errors[0].Message)
	}

	if info.Subtasks != nil && len(info.Subtasks) > 0 {
		if info.Subtasks[0].SubtaskID == "LoginEnterAlternateIdentifierSubtask" {
			err = fmt.Errorf("auth error: %v", "LoginEnterAlternateIdentifierSubtask")
		} else if info.Subtasks[0].SubtaskID == "LoginAcid" {
			err = fmt.Errorf("auth error: %v", "LoginAcid")
		}
	}

	return info.FlowToken, err
}

// IsLoggedIn check if scraper logged in
func (s *Scraper) IsLoggedIn() bool {
	s.isLogged = true
	s.setBearerToken(bearerToken2)
	req, err := http.NewRequest("GET", "https://api.twitter.com/1.1/account/verify_credentials.json", nil)
	if err != nil {
		return false
	}
	var verify verifyCredentials
	err = s.RequestAPI(req, &verify)
	if err != nil || verify.Errors != nil {
		s.isLogged = false
		s.setBearerToken(bearerToken)
	} else {
		s.isLogged = true
	}
	return s.isLogged
}

// Login to Twitter
// Use Login(username, password) for ordinary login
// or Login(username, password, email) for login if you have email confirmation
func (s *Scraper) Login(credentials ...string) error {
	var username, password, email string
	if len(credentials) == 2 {
		username = credentials[0]
		password = credentials[1]
	} else if len(credentials) == 3 {
		username = credentials[0]
		password = credentials[1]
		email = credentials[2]
	} else {
		return fmt.Errorf("invalid credentials")
	}

	s.setBearerToken(bearerToken2)

	err := s.GetGuestToken()
	if err != nil {
		return err
	}

	// flow start
	data := map[string]interface{}{
		"flow_name": "login",
		"input_flow_data": map[string]interface{}{
			"flow_context": map[string]interface{}{
				"debug_overrides": map[string]interface{}{},
				"start_location":  map[string]interface{}{"location": "splash_screen"},
			},
		},
	}
	flowToken, err := s.getFlowToken(data)
	if err != nil {
		return err
	}

	// flow instrumentation step
	data = map[string]interface{}{
		"flow_token": flowToken,
		"subtask_inputs": []map[string]interface{}{
			{
				"subtask_id":         "LoginJsInstrumentationSubtask",
				"js_instrumentation": map[string]interface{}{"response": "{}", "link": "next_link"},
			},
		},
	}
	flowToken, err = s.getFlowToken(data)
	if err != nil {
		return err
	}

	// flow username step
	data = map[string]interface{}{
		"flow_token": flowToken,
		"subtask_inputs": []map[string]interface{}{
			{
				"subtask_id": "LoginEnterUserIdentifierSSO",
				"settings_list": map[string]interface{}{
					"setting_responses": []map[string]interface{}{
						{
							"key":           "user_identifier",
							"response_data": map[string]interface{}{"text_data": map[string]interface{}{"result": username}},
						},
					},
					"link": "next_link",
				},
			},
		},
	}
	flowToken, err = s.getFlowToken(data)
	if err != nil {
		return err
	}

	// flow password step
	data = map[string]interface{}{
		"flow_token": flowToken,
		"subtask_inputs": []map[string]interface{}{
			{
				"subtask_id":     "LoginEnterPassword",
				"enter_password": map[string]interface{}{"password": password, "link": "next_link"},
			},
		},
	}
	flowToken, err = s.getFlowToken(data)
	if err != nil {
		return err
	}

	// flow duplication check
	data = map[string]interface{}{
		"flow_token": flowToken,
		"subtask_inputs": []map[string]interface{}{
			{
				"subtask_id":              "AccountDuplicationCheck",
				"check_logged_in_account": map[string]interface{}{"link": "AccountDuplicationCheck_false"},
			},
		},
	}
	flowToken, err = s.getFlowToken(data)
	if err != nil {
		if strings.Contains(err.Error(), "LoginAcid") {
			// flow acid
			data = map[string]interface{}{
				"flow_token": flowToken,
				"subtask_inputs": []map[string]interface{}{
					{
						"subtask_id": "LoginAcid",
						"enter_text": map[string]interface{}{"text": email, "link": "next_link"},
					},
				},
			}
			_, err = s.getFlowToken(data)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	s.isLogged = true
	return nil
}

// Logout is reset session
func (s *Scraper) Logout() error {
	req, err := http.NewRequest("POST", logoutURL, nil)
	if err != nil {
		return err
	}
	err = s.RequestAPI(req, nil)
	if err != nil {
		return err
	}

	s.isLogged = false
	s.guestToken = ""
	s.client.Jar, _ = cookiejar.New(nil)
	s.setBearerToken(bearerToken)
	return nil
}

func (s *Scraper) GetCookies() []*http.Cookie {
	var cookies []*http.Cookie
	for _, cookie := range s.client.Jar.Cookies(twURL) {
		if strings.Contains(cookie.Name, "guest") {
			continue
		}
		cookie.Domain = twURL.Host
		cookies = append(cookies, cookie)
	}
	return cookies
}

func (s *Scraper) SetCookies(cookies []*http.Cookie) {
	s.client.Jar.SetCookies(twURL, cookies)
}
