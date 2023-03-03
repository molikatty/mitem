package mitem

import (
	"sync/atomic"
	"unsafe"
)

type Int interface {
	int | int32 | int64
}

type Uint interface {
	uint | uint32 | uint64
}

// linear iteration
type Scanner[T any, P Int | Uint] struct {
	// min number
	min atomic.Uint64
	// max number
	max uint64
	// min item to max
	f func(P) T
	// index
	index P
	// data
	data T
}

func NewScan[T any, P Int | Uint](min, max P, f func(P) T) *Scanner[T, P] {
	return &Scanner[T, P]{
		max: uint64(max),
		min: *(*atomic.Uint64)(unsafe.Pointer(&min)),
		f:   f,
	}
}

// check if the next element exists
func (scan *Scanner[T, P]) Next() bool {
	index := scan.min.Add(1) - 1
	if index > scan.max {
		// avoid overflow caused by high concurrency
		// min = 18446744073709551615 + 1 == 0
		// min < max
		scan.min.Store(scan.max + 1)
		return false
	}
	scan.data = scan.f(P(index))
	scan.index = P(index)
	return true
}

// get data
func (scan *Scanner[T, P]) Data() (T, P) {
	return scan.data, scan.index
}
