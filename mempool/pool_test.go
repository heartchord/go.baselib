// +build !race

package mempool

import (
	"sync"
	"testing"
)

func testNewFunc() interface{} {
	return 1234567890
}

func Benchmark_Pool(b *testing.B) {
	p := NewPool(2, testNewFunc)
	var v interface{}
	v = 1
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Put(v)
			p.Get()
		}
	})
}

func BenchmarkBuiltinPool(b *testing.B) {
	var p sync.Pool
	var v interface{}
	v = 1
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Put(v)
			p.Get()
		}
	})
}
