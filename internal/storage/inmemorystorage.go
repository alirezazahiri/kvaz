package storage

import "fmt"

type InMemoryStorage struct {
	cache map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		cache: make(map[string]string),
	}
}

func (ims *InMemoryStorage) Get(k string) (string, error) {
	if v, ok := ims.cache[k]; ok {
		return v, nil
	}
	return "", fmt.Errorf("key '%s' not found", k)
}

func (ims *InMemoryStorage) Set(k, v string) bool {
	if _, ok := ims.cache[k]; !ok {
		ims.cache[k] = v
		return true
	}
	return false
}
