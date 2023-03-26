package memory

import (
	"remote-force/server/server"

	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
)

type store struct {
	c *cache.Cache
}

func (s *store) Save(id string, o *oauth2.Token) error {
	s.c.SetDefault(id, o)
	return nil
}

func New() server.Store {
	c := cache.New(cache.NoExpiration, cache.NoExpiration)

	return &store{
		c: c,
	}
}
