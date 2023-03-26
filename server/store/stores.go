package store

import (
	"remote-force/server/entity"

	"golang.org/x/oauth2"
)

type User Store[entity.User]
type Token Store[*oauth2.Token]
