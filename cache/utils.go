package cacheutils

import (
	"hash/fnv"

	intutils "github.com/sudosz/go-utils/ints"
)

func getKeyHash(key []byte) []byte {
	keyHash := fnv.New64()
	keyHash.Write(key)
	return intutils.Int64ToBytes(int64(keyHash.Sum64()))
}
