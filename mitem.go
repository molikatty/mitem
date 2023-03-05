package mitem

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/molikatty/spinlock"
)

type Int interface {
	int | int32 | int64
}

type Uint interface {
	uint | uint32 | uint64
}

type Integer interface {
	Int | Uint
}

type Item[T any, P Integer] interface {
	next(*T) bool
}

type Scan[T any, P Integer] struct {
	// Iterator
	Item[T, P]
	// Ensure that data is not stored in a goroutine
	// manner to avoid data race.
	data T
}

// linear iteration
type Iterator[T any, P Integer] struct {
	proto atomic.Uint64
	// min number
	min atomic.Uint64
	// max number
	max uint64
	// min item to max
	handle func(P) T
	// lock
	l sync.Locker
	// maximum of reset
	re uint64
	// count reset
	has atomic.Uint64
}

func NewScan[T any, P Integer](min, max P, handle func(P) T) *Iterator[T, P] {
	m := *(*atomic.Uint64)(unsafe.Pointer(&min))
	return &Iterator[T, P]{
		max:    uint64(max),
		proto:  m,
		min:    m,
		handle: handle,
		l:      spinlock.NewLock(),
	}
}

func (scan *Scan[T, P]) Next() bool {
	return scan.next(&scan.data)
}

func (scan *Scan[T, P]) Data() T {
	return scan.data
}

func (item *Iterator[T, P]) Item() *Scan[T, P] {
	return &Scan[T, P]{Item: item}
}

// Check if there is still another element.
func (item *Iterator[T, P]) next(data *T) bool {
	index := item.min.Add(1) - 1
	if index > item.max {
		// avoid overflow caused by high concurrency
		// min = 18446744073709551615 + 1 == 0
		// min < max
		item.min.Store(item.max + 1)
		return false
	}
	*data = item.handle(P(index))
	return true
}

func (item *Iterator[T, P]) SetResetNum(num P) {
	item.re = uint64(num)
}

// Reset the iterator
func (item *Iterator[T, P]) Reset() bool {
	item.l.Lock()
	defer item.l.Unlock()
	if item.min.Load() > item.max && item.has.Add(1)-1 < item.re {
		item.min = item.proto
		return true
	}

	return false
}
