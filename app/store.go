package main

import (
	"time"
)

type Store struct {
	storage map[string]V
}

type V struct {
	value    string
	expiring *time.Time
}

func NewStore() *Store {
	return &Store{
		storage: make(map[string]V),
	}
}

func (s *Store) Get(key string) string {
	v := s.storage[key]
	if v.expiring == nil || v.expiring.After(time.Now()) {
		return v.value
	} else {
		delete(s.storage, key)
	}

	return ""
}

func (s *Store) Set(key, value string, expiring time.Duration) {
	if expiring == time.Duration(0) {
		s.storage[key] = V{value: value}
		return
	} else {
		expiringTime := time.Now().Add(expiring)
		s.storage[key] = V{value: value, expiring: &expiringTime}
	}
}
