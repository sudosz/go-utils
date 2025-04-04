package channelutils

import (
	"reflect"
	"sync"
)

var (
	pools = make(map[string]any)
)

// getPoolOf returns a sync.Pool for the given type T, creating it if it doesn’t exist.
// Optimization: Uses a global map to reuse pools per type.
func getPoolOf[T comparable](t T) *sync.Pool {
	ts := reflect.TypeOf(t).String()
	if p, ok := pools[ts]; ok {
		return p.(*sync.Pool)
	} else {
		v := new(sync.Pool)
		pools[ts] = v
		return v
	}
}

// New creates a new channel of type T with optional buffering.
// Optimization: Single allocation for channel creation.
func New[T any](buffered ...int) *chan T {
	var ch chan T
	if len(buffered) > 0 {
		ch = make(chan T, buffered[0])
	} else {
		ch = make(chan T)
	}
	return &ch
}

// AcquireChannel retrieves a channel from the pool or creates a new one if none available.
// Optimization: Reuses channels via pooling to reduce allocations.
func AcquireChannel[T any]() *chan T {
	if ch := getPoolOf(new(chan T)).Get(); ch != nil {
		return ch.(*chan T)
	}
	return New[T]()
}

// ReleaseChannel returns a channel to the pool for reuse.
// Optimization: Enables channel reuse to minimize garbage collection.
func ReleaseChannel[T any](ch *chan T) {
	getPoolOf(ch).Put(ch)
}
