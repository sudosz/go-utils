package pool

import (
	"container/list"
	"time"
)

type queued struct {
	when  time.Time
	slice *[]byte
}

func makeRecycler(length int) chan *[]byte {

	ch := make(chan *[]byte)

	go func() {
		q := new(list.List)
		for {
			if q.Len() == 0 {
				b := make([]byte, length)
				q.PushFront(queued{when: time.Now(), slice: &b})
			}

			e := q.Front()

			timeout := time.NewTimer(time.Minute)
			select {
			case b := <-ch:
				timeout.Stop()
				q.PushFront(queued{when: time.Now(), slice: b})

			case ch <- e.Value.(queued).slice:
				timeout.Stop()
				q.Remove(e)

			case <-timeout.C:
				e := q.Front()
				for e != nil {
					n := e.Next()
					if time.Since(e.Value.(queued).when) > time.Minute {
						q.Remove(e)
						e.Value = nil
					}
					e = n
				}
			}
		}

	}()

	return ch
}

type RecyclerBufferPool struct {
	ch chan *[]byte
}

func NewBufferPool(length int) *RecyclerBufferPool {
	return &RecyclerBufferPool{
		ch: makeRecycler(length),
	}
}

func (p *RecyclerBufferPool) Get() []byte {
	return *<-p.ch
}

func (p *RecyclerBufferPool) Put(buf *[]byte) {
	p.ch <- buf
}
