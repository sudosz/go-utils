package poolutils

import (
	"sync"
	"sync/atomic"
)

const DefaultTypedLimitedPoolNumber = 1 << 7

type TypedLimitedPool[T any] struct {
	New func() *T
	N   int

	once sync.Once
	pool sync.Pool
	mux  sync.Mutex
	n    atomic.Uint32
}

func newF[T any](new func() *T) func() any {
	return func() any {
		return new()
	}
}

func NewTypedLimitedPool[T any](n int, new func() *T) *TypedLimitedPool[T] {
	return &TypedLimitedPool[T]{
		N:   n,
		New: new,
	}
}

func (p *TypedLimitedPool[T]) init() {
	p.once.Do(func() {
		p.pool = sync.Pool{
			New: newF(p.New),
		}
	})
}

func (p *TypedLimitedPool[T]) Get() *T {
	p.init()
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.n.Load() > 0 {
		p.n.Add(^uint32(0))
	}
	return p.pool.Get().(*T)
}

func (p *TypedLimitedPool[T]) Put(v *T) {
	p.init()
	p.once.Do(func() {
		p.pool = sync.Pool{
			New: newF(p.New),
		}
	})
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.n.Load() < uint32(p.N) {
		p.n.Add(1)
		p.pool.Put(v)
	}
}
