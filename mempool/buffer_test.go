package mempool

import "testing"

func Benchmark_Buffer(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := NewBuffer(defaultBytesPool, 8, 1024)
			buf.DecRef()
		}
	})
}
