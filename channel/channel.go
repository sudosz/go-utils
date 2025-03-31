package channel

import (
	"reflect"
	"sync"
)

var (
	pools = make(map[string]any)
)

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

func New[T any](buffered ...int) *chan T {
	var ch chan T
	if len(buffered) > 0 {
		ch = make(chan T, buffered[0])
	} else {
		ch = make(chan T)
	}
	return &ch
}

func AcquireChannel[T any]() *chan T {
	if ch := getPoolOf(new(chan T)).Get(); ch != nil {
		return ch.(*chan T)
	}
	return New[T]()
}

func ReleaseChannel[T any](ch *chan T) {
	getPoolOf(ch).Put(ch)
}
