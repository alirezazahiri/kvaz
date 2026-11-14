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

func (ims *InMemoryStorage) Delete(k string) bool {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	if _, ok := ims.cache[k]; ok {
		delete(ims.cache, k)
		return ok
	}
	return false
}

func (ims *InMemoryStorage) Exists(k string) bool {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	_, ok := ims.cache[k]
	return ok
}
