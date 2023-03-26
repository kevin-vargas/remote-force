package config

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

const (
	oauth_client_id       = "OAUTH_CLIENT_ID"
	oauth_auth_device_url = "OAUTH_AUTH_DEVICE_URL"
	oauth_scopes          = "OAUTH_SCOPES"
	addr                  = "ADDR"
	available_commands    = "AVAILABLE_COMMANDS"
	secret_key            = "SECRET_KEY"
)

var (
	oauth_endpoint = endpoints.GitHub
)

const (
	token = " "
)

type Oauth struct {
	CliendID      string
	AuthDeviceURL string
	Scopes        []string
	Endpoint      oauth2.Endpoint
}

type Config struct {
	Oauth             Oauth
	SecretKey         string
	ADDR              string
	AvailableCommands []string
}

func New() (cfg Config, err error) {
	ac, err := osGetEnvArray(available_commands, token)
	if err != nil {
		return
	}
	p, err := osGetEnv(addr)
	if err != nil {
		return
	}
	j, err := osGetEnv(secret_key)
	if err != nil {
		return
	}
	cid, err := osGetEnv(oauth_client_id)
	if err != nil {
		return
	}
	scs, err := osGetEnvArray(oauth_scopes, token)
	if err != nil {
		return
	}
	durl, err := osGetEnv(oauth_auth_device_url)
	if err != nil {
		return
	}
	return Config{
		Oauth: Oauth{
			Scopes:        scs,
			CliendID:      cid,
			AuthDeviceURL: durl,
			Endpoint:      oauth_endpoint,
		},
		SecretKey:         j,
		ADDR:              p,
		AvailableCommands: ac,
	}, nil
}

func osGetEnvArray(s string, token string) ([]string, error) {
	r, err := osGetEnv(s)
	if err != nil {
		return nil, err
	}
	return strings.Split(r, token), nil
}

func osGetEnv(s string) (string, error) {
	r := os.Getenv(s)
	if r == "" {
		return "", fmt.Errorf("missing env var %s", s)
	}
	return r, nil
}
