package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	addr               = "ADDR"
	available_commands = "AVAILABLE_COMMANDS"
)

const (
	token = " "
)

type Config struct {
	ADDR              string
	AvailableCommands []string
}

func New() (Config, error) {
	ac, err := osGetEnvArray(available_commands, token)
	if err != nil {
		return Config{}, err
	}
	p, err := osGetEnv(addr)
	if err != nil {
		return Config{}, err
	}
	return Config{
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
