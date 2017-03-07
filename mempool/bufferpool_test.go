package mempool

import "testing"

func Benchmark_BufferPool_32(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := GetFromDefaultBufferPool(8, 32)
			buf.DecRef()
		}
	})
}

func Benchmark_BufferPool_1024(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := GetFromDefaultBufferPool(8, 1024)
			buf.DecRef()
		}
	})
}

func Benchmark_BufferPool_LargeSize(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := GetFromDefaultBufferPool(8, 65536+1)
			buf.DecRef()
		}
	})
}
