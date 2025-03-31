package pool

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DefaultLRULimitedPoolNumber = 1 << 7
	DefaultLimitedPoolLRU       = 10 * time.Second
)

var emptyTime time.Time

type LRULimitedPool struct {
	N        int
	Interval time.Duration
	New      func() any
	Deleter  func(any)

	lru    time.Time
	once   sync.Once
	pool   sync.Pool
	mux    sync.Mutex
	lrumux sync.Mutex
	n      atomic.Uint32
}

func NewLRULimitedPool(n int, interval time.Duration, new func() any, deleter ...func(any)) *LRULimitedPool {
	df := (func(v any))(nil)
	if len(deleter) > 0 {
		df = deleter[0]
	}
	return &LRULimitedPool{
		N:        n,
		Interval: interval,
		New:      new,
		Deleter:  df,
	}
}

func (p *LRULimitedPool) init() {
	p.once.Do(func() {
		p.pool = sync.Pool{
			New: p.New,
		}
		if p.Interval > 0 {
			go func() {
				for {
					time.Sleep(p.Interval)
					p.cleanup()
				}
			}()
		}
	})
}

func (p *LRULimitedPool) cleanup() {
	p.lrumux.Lock()
	defer func() { p.lrumux.Unlock() }()
	if !p.lru.IsZero() && time.Since(p.lru) > p.Interval {
		p.Cleanup()
		p.lru = emptyTime
	}
}

func (p *LRULimitedPool) Get() any {
	p.init()
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.n.Load() > 0 {
		p.n.Add(^uint32(0))
	}
	return p.pool.Get()
}

func (p *LRULimitedPool) Put(v any) {
	p.init()
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.n.Load() < uint32(p.N) {
		p.n.Add(1)
		p.pool.Put(v)
		p.lrumux.Lock()
		p.lru = time.Now()
		p.lrumux.Unlock()
	}
}

func (p *LRULimitedPool) Cleanup() {
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.n.Load() > 0 {
		for i := 0; i < int(p.n.Load()); i++ {
			if p.Deleter == nil {
				p.pool.Get()
			} else {
				p.Deleter(p.pool.Get())
			}
		}
		p.n.Store(0)
	}
}

func NewLRULimitedBufferPool(n int, size int, interval time.Duration) *LRULimitedPool {
	return &LRULimitedPool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, size))
		},
		N:        n,
		Interval: interval,
	}
}
