package memory

import (
	"errors"
	"remote-force/server/store"

	"github.com/patrickmn/go-cache"
)

type implementation[T any] struct {
	c *cache.Cache
}

func (s *implementation[T]) Save(id string, v T) error {
	s.c.SetDefault(id, v)
	return nil
}

func (s *implementation[T]) Get(id string) (T, bool, error) {
	var t T
	v, found := s.c.Get(id)
	if !found {
		return t, false, nil
	}
	if t, ok := v.(T); ok {
		return t, true, nil
	}
	return t, false, errors.New("invalid type from cache")
}

func New[T any]() store.Store[T] {
	c := cache.New(cache.NoExpiration, cache.NoExpiration)

	return &implementation[T]{
		c: c,
	}
}
