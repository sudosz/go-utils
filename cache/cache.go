package cache

import (
	"time"

	"git.mills.io/prologic/bitcask"
	"github.com/sudosz/go-utils/bytes"
)

// Cache represents a cache using bitcask as the underlying storage.
type Cache struct {
	DB *bitcask.Bitcask
}

// New creates a new Cache instance with the given folder path.
// Optimization: Relies on bitcask's efficiency; no additional overhead added.
func New(folder string) (*Cache, error) {
	db, err := bitcask.Open(folder)
	if err != nil {
		return nil, err
	}
	return &Cache{DB: db}, nil
}

// Get retrieves the value for the given string key, converting it to bytes.
// Optimization: Uses zero-copy S2b for key conversion.
func (c *Cache) Get(key string) ([]byte, error) {
	return c.GetBytes(bytes.S2b(key))
}

// GetBytes retrieves the value for the given byte slice key using a hash.
// Optimization: Efficient key hashing via getKeyHash.
func (c *Cache) GetBytes(key []byte) ([]byte, error) {
	return c.DB.Get(getKeyHash(key))
}

// Has checks if the given string key exists in the cache.
// Optimization: Uses zero-copy S2b for key conversion.
func (c *Cache) Has(key string) bool {
	return c.HasBytes(bytes.S2b(key))
}

// HasBytes checks if the given byte slice key exists in the cache.
// Optimization: Efficient key hashing via getKeyHash.
func (c *Cache) HasBytes(key []byte) bool {
	return c.DB.Has(getKeyHash(key))
}

// Set sets the value for the given string key and value.
// Optimization: Uses zero-copy S2b for both key and value.
func (c *Cache) Set(key string, value string) error {
	return c.SetBytesKV(bytes.S2b(key), bytes.S2b(value))
}

// SetBytesKVWithTTL sets the value for the byte slice key with a time-to-live.
// Optimization: Direct use of bitcask’s TTL feature.
func (c *Cache) SetBytesKVWithTTL(key []byte, value []byte, ttl time.Duration) error {
	return c.DB.PutWithTTL(getKeyHash(key), value, ttl)
}

// SetBytesK sets the value for a byte slice key with a string value.
// Optimization: Zero-copy conversion for value.
func (c *Cache) SetBytesK(key []byte, value string) error {
	return c.SetBytesKV(key, bytes.S2b(value))
}

// SetBytesV sets the value for a string key with a byte slice value.
// Optimization: Zero-copy conversion for key.
func (c *Cache) SetBytesV(key string, value []byte) error {
	return c.SetBytesKV(bytes.S2b(key), value)
}

// SetBytesKV sets the value for a byte slice key and value.
// Optimization: Efficient key hashing via getKeyHash.
func (c *Cache) SetBytesKV(key []byte, value []byte) error {
	return c.DB.Put(getKeyHash(key), value)
}

// Close closes the cache database.
// Optimization: Direct passthrough to bitcask.Close.
func (c *Cache) Close() error {
	return c.DB.Close()
}

// DelAll deletes all keys in the cache.
// Optimization: Relies on bitcask’s efficient deletion.
func (c *Cache) DelAll() error {
	return c.DB.DeleteAll()
}

// Del deletes the given string key from the cache.
// Optimization: Zero-copy conversion for key.
func (c *Cache) Del(key string) error {
	return c.DelBytes(bytes.S2b(key))
}

// DelBytes deletes the given byte slice key from the cache.
// Optimization: Efficient key hashing via getKeyHash.
func (c *Cache) DelBytes(key []byte) error {
	return c.DB.Delete(key)
}

// Len returns the number of keys in the cache.
// Optimization: Direct passthrough to bitcask.Len.
func (c *Cache) Len() int {
	return c.DB.Len()
}

// RunGC runs the garbage collector on the cache.
// Optimization: Direct passthrough to bitcask.RunGC.
func (c *Cache) RunGC() error {
	return c.DB.RunGC()
}
