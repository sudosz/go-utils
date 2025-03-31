package cache

import (
	"hash/fnv"

	"github.com/sudosz/go-utils/ints"
)

func getKeyHash(key []byte) []byte {
	keyHash := fnv.New64()
	keyHash.Write(key)
	return ints.Int64ToBytes(int64(keyHash.Sum64()))
}
