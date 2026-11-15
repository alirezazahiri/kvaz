package storage_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"kvaz/internal/storage"
)

func prepareStorage() storage.StorageEngine {
	return storage.NewInMemoryStorage()
}

func TestNewInMemoryStorage(t *testing.T) {
	store := prepareStorage()
	if store == nil {
		t.Fatal("expected non-nil storage")
	}
	if _, ok := store.Get("anything"); ok {
		t.Fatal("expected empty storage to miss keys")
	}
}

func TestSetReturnsTrue(t *testing.T) {
	store := prepareStorage()
	if !store.Set("k", "v") {
		t.Fatal("set should return true")
	}
}

func TestSetAndGet(t *testing.T) {
	store := prepareStorage()

	store.Set("name", "Alireza")

	got, ok := store.Get("name")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if got != "Alireza" {
		t.Fatalf("expected %q, got %q", "Alireza", got)
	}
}

func TestSetOverwritesExistingValue(t *testing.T) {
	store := prepareStorage()

	store.Set("name", "first")
	store.Set("name", "second")

	got, ok := store.Get("name")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if got != "second" {
		t.Fatalf("expected %q, got %q", "second", got)
	}
}

func TestGetMissingKey(t *testing.T) {
	store := prepareStorage()

	got, ok := store.Get("missing")
	if ok {
		t.Fatal("expected missing key")
	}
	if got != "" {
		t.Fatalf("expected empty value, got %q", got)
	}
}

func TestDeleteExistingKey(t *testing.T) {
	store := prepareStorage()
	store.Set("k", "v")

	if !store.Delete("k") {
		t.Fatal("delete should return true for existing key")
	}
	if _, ok := store.Get("k"); ok {
		t.Fatal("key should be gone after Delete")
	}
	if store.Exists("k") {
		t.Fatal("exists should be false after Delete")
	}
}

func TestDeleteMissingKey(t *testing.T) {
	store := prepareStorage()
	if store.Delete("missing") {
		t.Fatal("Delete should return false for missing key")
	}
}

func TestDeleteAlreadyExpiredKey(t *testing.T) {
	store := prepareStorage()
	store.Set("k", "v")
	store.Expire("k", -1)
	_, state := store.TTL("k")
	if state != storage.RecordStateMissing {
		t.Fatalf("expected missing after expiry, got state %v", state)
	}
	if store.Delete("k") {
		t.Fatal("Delete should return false when TTL reports missing")
	}
}

func TestExistsPersistentKey(t *testing.T) {
	store := prepareStorage()
	store.Set("k", "v")
	if !store.Exists("k") {
		t.Fatal("Exists should be true for set key without deadline")
	}
}

func TestExistsMissingKey(t *testing.T) {
	store := prepareStorage()
	if store.Exists("missing") {
		t.Fatal("Exists should be false for missing key")
	}
}

func TestExistsAfterTTLExpiry(t *testing.T) {
	store := prepareStorage()
	store.Set("k", "v")
	store.Expire("k", -1)
	store.TTL("k")
	if store.Exists("k") {
		t.Fatal("Exists should be false after record expires")
	}
}

func TestExpireOnExistingKey(t *testing.T) {
	store := prepareStorage()
	store.Set("k", "v")

	if !store.Expire("k", 3600) {
		t.Fatal("Expire should return true for existing key")
	}

	ttl, state := store.TTL("k")
	if state != storage.RecordStateExpiring {
		t.Fatalf("expected RecordStateExpiring, got %v", state)
	}
	if ttl <= 0 || ttl > 3600 {
		t.Fatalf("expected TTL in (0, 3600], got %d", ttl)
	}

	got, ok := store.Get("k")
	if !ok || got != "v" {
		t.Fatalf("value should remain after Expire, got %q ok=%v", got, ok)
	}
}

func TestExpireMissingKey(t *testing.T) {
	store := prepareStorage()
	if store.Expire("missing", 10) {
		t.Fatal("Expire should return false for missing key")
	}
}

func TestExpireWhenTTLReportsMissing(t *testing.T) {
	store := prepareStorage()
	store.Set("k", "v")
	store.Expire("k", -1)
	store.TTL("k")
	if store.Expire("k", 10) {
		t.Fatal("Expire should return false when key was invalidated by TTL")
	}
}

func TestTTLMissingKey(t *testing.T) {
	store := prepareStorage()
	ttl, state := store.TTL("missing")
	if state != storage.RecordStateMissing {
		t.Fatalf("expected RecordStateMissing, got %v", state)
	}
	if ttl != 0 {
		t.Fatalf("expected TTL 0, got %d", ttl)
	}
}

func TestTTLPersistentKey(t *testing.T) {
	store := prepareStorage()
	store.Set("k", "v")

	ttl, state := store.TTL("k")
	if state != storage.RecordStatePersistent {
		t.Fatalf("expected RecordStatePersistent, got %v", state)
	}
	if ttl != 0 {
		t.Fatalf("expected TTL 0 for persistent key, got %d", ttl)
	}
}

func TestTTLExpiringKey(t *testing.T) {
	store := prepareStorage()
	store.Set("k", "v")
	store.Expire("k", 30)

	ttl, state := store.TTL("k")
	if state != storage.RecordStateExpiring {
		t.Fatalf("expected RecordStateExpiring, got %v", state)
	}
	if ttl <= 0 || ttl > 30 {
		t.Fatalf("expected TTL in (0, 30], got %d", ttl)
	}
}

func TestTTLRemovesExpiredKey(t *testing.T) {
	store := prepareStorage()
	store.Set("k", "v")
	store.Expire("k", -1)

	ttl, state := store.TTL("k")
	if state != storage.RecordStateMissing || ttl != 0 {
		t.Fatalf("expected missing TTL after deadline passed, got ttl=%d state=%v", ttl, state)
	}
	if _, ok := store.Get("k"); ok {
		t.Fatal("expired key should be removed from cache")
	}
}

func TestGetAfterNaturalExpiry(t *testing.T) {
	store := prepareStorage()
	store.Set("k", "v")
	store.Expire("k", 1)
	time.Sleep(1100 * time.Millisecond)

	got, ok := store.Get("k")
	if ok {
		t.Fatalf("expected expired key to be missing, got %q", got)
	}
}

func TestConcurrentAccess(t *testing.T) {
	store := prepareStorage()
	testfail := make(chan string)

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		select {
		case msg := <-testfail:
			t.Fatal(msg)
		default:
			break
		}
		go func(i int) {
			defer wg.Done()
			key := "shared"
			value := fmt.Sprintf("value%d", i)
			store.Set(key, value)
			v, _ := store.Get(key)
			if v != value {
				testfail <- fmt.Sprintf("expected value to be %s, but got %s", value, v)
			}
		}(i)
	}
	wg.Wait()

	_, ok := store.Get("shared")
	if !ok {
		t.Fatal("expected shared key to exist after concurrent writes")
	}
}

func TestConcurrentTTLVsMutations(t *testing.T) {
	store := prepareStorage()
	const workers = 8
	const iterations = 5000

	var wg sync.WaitGroup
	wg.Add(workers * 3)
	for range workers {
		go func() {
			defer wg.Done()
			for i := range iterations {
				key := fmt.Sprintf("k%d", i%50)
				store.Set(key, "v")
				store.Expire(key, 3600)
			}
		}()
		go func() {
			defer wg.Done()
			for i := range iterations {
				key := fmt.Sprintf("k%d", i%50)
				store.Delete(key)
			}
		}()
		go func() {
			defer wg.Done()
			for i := range iterations {
				key := fmt.Sprintf("k%d", i%50)
				_, _ = store.Get(key)
				_ = store.Delete(key)
				_ = store.Expire(key, 60)
			}
		}()
	}
	wg.Wait()
}
