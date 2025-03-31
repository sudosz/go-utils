package iter

import "sync"

type CycleIterator[T any] struct {
	mux    *sync.RWMutex
	values []T
	index  int
}

func NewCycleIterator[T any](values ...T) *CycleIterator[T] {
	return &CycleIterator[T]{
		values: values,
		mux:    &sync.RWMutex{},
	}
}

func (it *CycleIterator[T]) Add(values ...T) {
	it.mux.Lock()
	defer it.mux.Unlock()
	it.values = append(it.values, values...)
}

func (it *CycleIterator[T]) Get() (value T) {
	if len(it.values) == 0 {
		it.index = 0
		return value
	}

	if it.index >= len(it.values)-1 {
		it.index = 0
	} else {
		it.index++
	}
	return it.values[it.index]
}

func (it *CycleIterator[T]) Len() int {
	return len(it.values)
}
