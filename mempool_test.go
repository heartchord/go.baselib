package goblazer

import (
	"math/rand"
	"testing"
)

func Test_MemoryPool(t *testing.T) {
	mp := NewMemoryPool()

	mb, ok := mp.Allocate(4)
	if ok {
		mp.Recycle(mb)
	}
}

func Benchmark_MemoryPool(b *testing.B) {
	b.StopTimer()

	mp := NewMemoryPool()
	var mb *MemoryBlock

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		mb, _ = mp.Allocate(rand.Intn(64 * 1024))
		mp.Recycle(mb)
	}
}
