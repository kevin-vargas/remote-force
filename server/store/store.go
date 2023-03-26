package store

type Store[T any] interface {
	Save(key string, v T) error
	Get(key string) (T, bool, error)
}
