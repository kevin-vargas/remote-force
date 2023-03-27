package jwt

type Store interface {
	Save(jwt string) error
	Get() (string, bool, error)
	CleanUP() error
}
