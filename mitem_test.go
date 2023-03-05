package mitem

import (
	"math/rand"
	"sync"
	// "runtime"
	"sync/atomic"
	"testing"
	"time"
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
)

const (
	n = 1000
	t = 100
)

var (
	num       = randstr()
	rangeRune = []rune{}
	mitmByte  = []byte{}
	curMem    uint64
)

func randstr() []uint64 {
	rand.Seed(time.Now().Unix())
	var in []uint64
	for i := 0; i < 1e5; i++ {
		in = append(in, rand.Uint64())
	}
	return in
}

func BenchmarkRange(b *testing.B) {
	for j := 0; j < b.N; j++ {
		var k uint64
		for i := range num {
			k = num[i]
		}
		_ = k
	}
	// mem := runtime.MemStats{}
	// runtime.ReadMemStats(&mem)
	// curMem = mem.TotalAlloc/MiB - curMem
	// t.Logf("memory usage:%d MB", curMem)
}

func TestMitem(t *testing.T) {
	scan := NewScan(1, 1000, func(index uint64) uint64 {
		// return num[index]
		return index
	})

	scan.SetResetNum(1e3)
	var nm atomic.Uint64
	var wg sync.WaitGroup
	for j := 0; j < 1e3; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			item := scan.Item()
		again:
			for item.Next() {
				nm.Add(item.Data())
			}

			if scan.Reset() {
				goto again
			}
		}()
	}
	wg.Wait()
	t.Logf("nm: %d", nm.Load())
}
