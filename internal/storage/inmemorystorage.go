package storage

import (
	"fmt"
	"sync"
	"time"
)

type Record struct {
	Value    string
	Deadline time.Time
}
type InMemoryStorage struct {
	cache map[string]Record
	mu    sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		mu:    sync.RWMutex{},
		cache: make(map[string]Record),
	}
}

func (ims *InMemoryStorage) Get(k string) (string, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	if v, ok := ims.cache[k]; ok {
		return v.Value, nil
	}
	return "", fmt.Errorf("key '%s' not found", k)
}

func (ims *InMemoryStorage) Set(k, v string) bool {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.cache[k] = Record{Value: v}
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

func (ims *InMemoryStorage) Expire(k string, secs int64) bool {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	if v, ok := ims.cache[k]; ok {
		ims.cache[k] = Record{
			Value:    v.Value,
			Deadline: time.Now().Add(time.Duration(secs * int64(time.Second))),
		}
		return true
	}
	return false
}

func (ims *InMemoryStorage) TTL(k string) (int64, RecordState) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	// TODO: implement me!
	return 0, RecordStatePersistent
}
