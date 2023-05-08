package main

type Store struct {
	storage map[string]string
}

func NewStore() *Store {
	return &Store{
		storage: make(map[string]string),
	}
}

func (s *Store) Get(key string) string {
	return s.storage[key]
}

func (s *Store) Set(key, value string) {
	s.storage[key] = value
}
