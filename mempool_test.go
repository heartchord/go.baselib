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
		n := rand.Intn(64 * 1024)
		mb, _ = mp.Allocate(n)
		mp.Recycle(mb)
	}

	//mp.Statistics()
}
