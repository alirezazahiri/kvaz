package storage

import (
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

func (ims *InMemoryStorage) Get(k string) (string, bool) {
	if _, state := ims.TTL(k); state == RecordStateMissing {
		return "", false
	}
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	if v, ok := ims.cache[k]; ok {
		return v.Value, true
	}
	return "", false
}

func (ims *InMemoryStorage) Set(k, v string) bool {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.cache[k] = Record{Value: v}
	return true
}

func (ims *InMemoryStorage) Delete(k string) bool {
	if _, state := ims.TTL(k); state == RecordStateMissing {
		return false
	}
	ims.mu.Lock()
	defer ims.mu.Unlock()
	if _, ok := ims.cache[k]; ok {
		ims.invalidateKey(k)
		return ok
	}
	return false
}

func (ims *InMemoryStorage) Exists(k string) bool {
	if _, state := ims.TTL(k); state == RecordStateMissing {
		return false
	}
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	_, ok := ims.cache[k]
	return ok
}

func (ims *InMemoryStorage) Expire(k string, secs int64) bool {
	if _, state := ims.TTL(k); state == RecordStateMissing {
		return false
	}
	ims.mu.Lock()
	defer ims.mu.Unlock()
	if v, ok := ims.cache[k]; ok {
		ddline := time.Now().Add(time.Duration(secs * int64(time.Second)))
		ims.cache[k] = Record{
			Value:    v.Value,
			Deadline: ddline,
		}
		return true
	}
	return false
}

func (ims *InMemoryStorage) TTL(k string) (int64, RecordState) {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	v, ok := ims.cache[k]
	if !ok {
		return 0, RecordStateMissing
	}
	if v.Deadline.IsZero() {
		return 0, RecordStatePersistent
	}
	remainingSeconds := time.Until(v.Deadline).Seconds()
	if remainingSeconds <= 0 {
		ims.invalidateKey(k)
		return 0, RecordStateMissing
	}
	return int64(remainingSeconds), RecordStateExpiring
}

// this one assumes that the caller function already acquires the lock and that the key already exists
func (ims *InMemoryStorage) invalidateKey(k string) {
	delete(ims.cache, k)
}
