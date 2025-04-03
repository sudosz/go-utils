package iter

import "sync"

// CycleIterator is a thread-safe iterator that cycles through a list of values.
type CycleIterator[T any] struct {
	mux    *sync.RWMutex
	values []T
	index  int
}

// NewCycleIterator creates a new iterator with the provided values.
// Optimization: Single allocation for values slice.
func NewCycleIterator[T any](values ...T) *CycleIterator[T] {
	return &CycleIterator[T]{
		values: values,
		mux:    &sync.RWMutex{},
	}
}

// Add appends values to the iterator’s list, thread-safely.
// Optimization: Mutex ensures safe concurrent access.
func (it *CycleIterator[T]) Add(values ...T) {
	it.mux.Lock()
	defer it.mux.Unlock()
	it.values = append(it.values, values...)
}

// Get returns the next value in the cycle, wrapping around to the start if needed.
// Optimization: Thread-safe with minimal locking overhead; could use atomic index but complexity increases.
func (it *CycleIterator[T]) Get() (value T) {
	it.mux.Lock()
	defer it.mux.Unlock()
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

// Len returns the number of values in the iterator, thread-safely.
// Optimization: Uses read lock for concurrent read safety.
func (it *CycleIterator[T]) Len() int {
	it.mux.RLock()
	defer it.mux.RUnlock()
	return len(it.values)
}
