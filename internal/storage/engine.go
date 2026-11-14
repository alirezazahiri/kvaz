// Package storage: this one is responsible for
// actually putting the value somewhere so that it could be accessed later
package storage

type StorageEngine interface {
	// Set creates the key if it doesn't exist and overwrites the key if it already exists,
	// returns `true` if it wrote something, and `false` if not
	Set(k, v string) bool
	// Get retrieves the stored value by a given key, returns a pair of `(value, error)`,
	// if the key exists, the value will be returned (if no internal error occurs in the way)
	// if the key is non-existent, then the value will be `nil` and an error would be returned
	Get(k string) (string, error)
	// Delete deletes a stored value at a given key, if the value exists it will return `true` after deletion
	// (if nothing else goes wrong in the way), if the value does not exist, it will return `false`
	Delete(k string) bool
	// Exists returns `true` if a key exists, and `false` if not
	Exists(k string) bool
}
