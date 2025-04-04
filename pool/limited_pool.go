package poolutils

import (
	"sync"
	"sync/atomic"
)

const DefaultLimitedPoolNumber = 1 << 7

// LimitedPool manages a pool with a fixed maximum number of objects.
type LimitedPool struct {
	N    int
	New  func() any
	once sync.Once
	pool sync.Pool
	mux  sync.Mutex
	n    atomic.Uint32
}

// NewLimitedPool creates a new LimitedPool with specified size and factory function.
// Optimization: Lazy initialization via sync.Once.
func NewLimitedPool(n int, new func() any) *LimitedPool {
	return &LimitedPool{
		N:   n,
		New: new,
	}
}

// init initializes the pool lazily on first use.
// Optimization: Ensures single initialization with minimal overhead.
func (p *LimitedPool) init() {
	p.once.Do(func() {
		p.pool = sync.Pool{New: p.New}
	})
}

// Get retrieves an object from the pool or creates a new one if empty.
// Optimization: Atomic counter ensures thread-safe object counting.
func (p *LimitedPool) Get() any {
	p.init()
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.n.Load() > 0 {
		p.n.Add(^uint32(0))
	}
	return p.pool.Get()
}

// Put returns an object to the pool if it’s not full.
// Optimization: Atomic check prevents overfilling.
func (p *LimitedPool) Put(v any) {
	p.init()
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.n.Load() < uint32(p.N) {
		p.n.Add(1)
		p.pool.Put(v)
	}
}
