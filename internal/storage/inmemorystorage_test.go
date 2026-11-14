package storage_test

import (
	"fmt"
	"kvaz/internal/storage"
	"sync"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	store := storage.NewInMemoryStorage()

	store.Set("name", "Alireza")

	got, err := store.Get("name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "Alireza" {
		t.Fatalf("expected %q, got %q", "Alireza", got)
	}
}

func TestSetOverwritesExistingValue(t *testing.T) {
	store := storage.NewInMemoryStorage()

	store.Set("name", "first")
	store.Set("name", "second")

	got, err := store.Get("name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "second" {
		t.Fatalf("expected %q, got %q", "second", got)
	}
}

func TestGetMissingKey(t *testing.T) {
	store := storage.NewInMemoryStorage()

	_, err := store.Get("missing")
	if err == nil {
		t.Fatal("expected an error for a missing key")
	}
}

func TestConcurrentAccess(t *testing.T) {
	store := storage.NewInMemoryStorage()

	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := "shared"
			store.Set(key, fmt.Sprintf("value%d", i))

			_, _ = store.Get(key)
		}(i)
	}

	fmt.Println(store.Get("shared"))

	wg.Wait()
}
