package cache

import (
	"time"

	"github.com/sudosz/go-utils/bytes"
	"git.mills.io/prologic/bitcask"
)

type Cache struct {
	DB *bitcask.Bitcask
}

func New(folder string) (*Cache, error) {
	db, err := bitcask.Open(folder)
	if err != nil {
		return nil, err
	}

	return &Cache{
		DB: db,
	}, nil
}

func (c *Cache) Get(key string) ([]byte, error) {
	return c.GetBytes(bytes.S2b(key))
}

func (c *Cache) GetBytes(key []byte) ([]byte, error) {
	return c.DB.Get(getKeyHash(key))
}

func (c *Cache) Has(key string) bool {
	return c.HasBytes(bytes.S2b(key))
}

func (c *Cache) HasBytes(key []byte) bool {
	return c.DB.Has(getKeyHash(key))
}

func (c *Cache) Set(key string, value string) error {
	return c.SetBytesKV(bytes.S2b(key), bytes.S2b(value))
}

func (c *Cache) SetBytesKVWithTTL(key []byte, value []byte, ttl time.Duration) error {
	return c.DB.PutWithTTL(getKeyHash(key), value, ttl)
}

func (c *Cache) SetBytesK(key []byte, value string) error {
	return c.SetBytesKV(key, bytes.S2b(value))
}

func (c *Cache) SetBytesV(key string, value []byte) error {
	return c.SetBytesKV(bytes.S2b(key), value)
}

func (c *Cache) SetBytesKV(key []byte, value []byte) error {
	return c.DB.Put(getKeyHash(key), value)
}

func (c *Cache) Close() error {
	return c.DB.Close()
}

func (c *Cache) DelAll() error {
	return c.DB.DeleteAll()
}

func (c *Cache) Del(key string) error {
	return c.DelBytes(bytes.S2b(key))
}

func (c *Cache) DelBytes(key []byte) error {
	return c.DB.Delete(key)
}

func (c *Cache) Len() int {
	return c.DB.Len()
}

func (c *Cache) RunGC() error {
	return c.DB.RunGC()
}
