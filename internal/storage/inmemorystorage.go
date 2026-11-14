package storage

import (
	"fmt"
	"sync"
)

type InMemoryStorage struct {
	cache map[string]string
	mu    sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		mu:    sync.RWMutex{},
		cache: make(map[string]string),
	}
}

func (ims *InMemoryStorage) Get(k string) (string, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	if v, ok := ims.cache[k]; ok {
		return v, nil
	}
	return "", fmt.Errorf("key '%s' not found", k)
}

func (ims *InMemoryStorage) Set(k, v string) bool {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.cache[k] = v
	return true
}
