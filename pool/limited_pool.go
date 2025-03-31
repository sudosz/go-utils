package pool

import (
	"sync"
	"sync/atomic"
)

const DefaultLimitedPoolNumber = 1 << 7

type LimitedPool struct {
	N   int
	New func() any

	once sync.Once
	pool sync.Pool
	mux  sync.Mutex
	n    atomic.Uint32
}

func NewLimitedPool(n int, new func() any) *LimitedPool {
	return &LimitedPool{
		N:   n,
		New: new,
	}
}

func (p *LimitedPool) init() {
	p.once.Do(func() {
		p.pool = sync.Pool{
			New: p.New,
		}
	})
}

func (p *LimitedPool) Get() any {
	p.init()
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.n.Load() > 0 {
		p.n.Add(^uint32(0))
	}
	return p.pool.Get()
}

func (p *LimitedPool) Put(v any) {
	p.init()
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.n.Load() < uint32(p.N) {
		p.n.Add(1)
		p.pool.Put(v)
	}
}
