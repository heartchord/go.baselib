package goblazer

import "testing"

/*
func Test_MemoryPool(t *testing.T) {
	mp := NewMemoryPool()

	mb, ok := mp.Allocate(4)
	if ok {
		mp.Recycle(mb)
	}
}
*/

var testmp = NewMemoryPool()

func Benchmark_MemoryPool_2014(b *testing.B) {
	var mb *MemoryBlock
	for i := 0; i < b.N; i++ {
		mb, _ = testmp.Allocate(1024)
		if mb != nil {
		}
		testmp.Recycle(mb)
	}
}

func Benchmark_MemoryPool_32(b *testing.B) {
	var mb *MemoryBlock
	for i := 0; i < b.N; i++ {
		mb, _ = testmp.Allocate(32)
		if mb != nil {
		}
		testmp.Recycle(mb)
	}
}

func Benchmark_BuiltinAlloc_1024(b *testing.B) {
	var mb []byte
	for i := 0; i < b.N; i++ {
		mb = make([]byte, 1024)
		if len(mb) != 0 {
		}
	}
}

func Benchmark_BuiltinAlloc_32(b *testing.B) {
	var mb []byte
	for i := 0; i < b.N; i++ {
		mb = make([]byte, 32)
		if len(mb) != 0 {
		}
	}
}
