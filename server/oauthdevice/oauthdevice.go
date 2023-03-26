package oauthdevice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"remote-force/server/config"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

var (
	errAccessDenied = errors.New("access denied by user")
	errTokenExpired = errors.New("token is expired")
)

type DeviceAuth struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete,omitempty"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval,omitempty"`
}

type Config struct {
	oc            oauth2.Config
	client        *http.Client
	DeviceAuthURL string
}

func New(cfg config.Oauth) *Config {
	return &Config{
		DeviceAuthURL: cfg.AuthDeviceURL,
		client:        http.DefaultClient,
		oc: oauth2.Config{
			ClientID: cfg.CliendID,
			Endpoint: cfg.Endpoint,
			Scopes:   cfg.Scopes,
		},
	}
}

func (c *Config) AuthDevice(ctx context.Context) (*DeviceAuth, error) {
	v := url.Values{
		"client_id": {c.oc.ClientID},
	}
	if len(c.oc.Scopes) > 0 {
		v.Set("scope", strings.Join(c.oc.Scopes, " "))
	}
	req, err := http.NewRequest(http.MethodPost, c.DeviceAuthURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"request for device code authorisation returned status %v (%v)",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	// Unmarshal response
	var dcr DeviceAuth
	if err := json.NewDecoder(resp.Body).Decode(&dcr); err != nil {

		return nil, err
	}
	return &dcr, nil
}

type tokenOrError struct {
	*oauth2.Token
	Error string `json:"error,omitempty"`
}

func (c *Config) Poll(ctx context.Context, da *DeviceAuth) (*oauth2.Token, error) {
	v := url.Values{
		"client_id":   {c.oc.ClientID},
		"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		"device_code": {da.DeviceCode},
	}
	if len(c.oc.Scopes) > 0 {
		v.Set("scope", strings.Join(c.oc.Scopes, " "))
	}
	// If no interval was provided, the client MUST use a reasonable default polling interval.
	// See https://tools.ietf.org/html/draft-ietf-oauth-device-flow-07#section-3.5
	interval := da.Interval
	if interval == 0 {
		interval = 5
	}

	for {
		time.Sleep(time.Duration(interval) * time.Second)
		req, err := http.NewRequest(http.MethodPost, c.oc.Endpoint.TokenURL, strings.NewReader(v.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")
		resp, err := c.client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("http error %v (%v) when polling for OAuth token",
				resp.StatusCode, http.StatusText(resp.StatusCode))
		}
		var token tokenOrError
		if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
			return nil, err
		}

		switch token.Error {
		case "":

			return token.Token, nil
		case "authorization_pending":

		case "slow_down":

			da.Interval *= 2
		case "access_denied":

			return nil, errAccessDenied
		case "expired_token":
			return nil, errTokenExpired
		default:

			return nil, fmt.Errorf("authorization failed: %v", token.Error)
		}
	}
}
