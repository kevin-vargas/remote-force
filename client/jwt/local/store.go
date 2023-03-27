package local

import (
	"os"
	"remote-force/client/jwt"
	"sync"
)

type store struct {
	sync.Mutex
	path string
}

func (s *store) Save(jwt string) error {
	s.Lock()
	defer s.Unlock()
	err := os.WriteFile(s.path, []byte(jwt), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) Get() (string, bool, error) {
	s.Lock()
	defer s.Unlock()
	if ok := fileExists(s.path); !ok {
		return "", false, nil
	}
	content, err := os.ReadFile(s.path)
	if err != nil {
		return "", false, err
	}
	return string(content), true, nil
}

func (s *store) CleanUP() error {
	return os.Remove(s.path)
}

func New(path string) jwt.Store {
	return &store{
		path: path,
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
