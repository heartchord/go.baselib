package goblazer

import (
	"math/rand"
	"testing"
)

func Test_MemoryPool(t *testing.T) {
	mp := NewMemoryPool()
	mp.InitPool()

	mb, ok := mp.Allocate(4)
	if ok {
		mp.Recycle(mb)
	}
}

func Benchmark_MemoryPool(b *testing.B) {
	b.StopTimer()

	mp := NewMemoryPool()
	mp.InitPool()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		mb, _ := mp.Allocate(rand.Intn(64 * 1024))
		//mp.Allocate(rand.Intn(64 * 1024))
		mp.Recycle(mb)
	}

	b.StopTimer()
	//mp.Statistics()

	//var mem runtime.MemStats
	//runtime.ReadMemStats(&mem)
	//fmt.Println(mem.Alloc)
	//fmt.Println(mem.TotalAlloc)
	//fmt.Println(mem.HeapAlloc)
	//fmt.Println(mem.HeapSys)
}
