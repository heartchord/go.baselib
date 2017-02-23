package mempool

import (
	"testing"
)

/*func Benchmark_BytesPool_0(b *testing.B) {
	b.StopTimer()
	bp := newBytesBlockList(1024)
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bb := bp.getFront()
			bp.putBack(bb)
		}
	})
}*/

func Benchmark_BytesPool_1(b *testing.B) {
	b.StopTimer()
	bp := newLocalBytesPool()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		bb := bp.getPrivate(1024)
		bp.putPrivate(bb)
	}
}

func Benchmark_BytesPool(b *testing.B) {
	b.StopTimer()
	bp := NewBytesPool()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		bb := bp.Get(1024)
		if bb != nil {
		}
		bp.Put(bb)
	}
}

func Benchmark_BuiltinMemPool(b *testing.B) {
	b.StopTimer()
	bp := NewBytesPool()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bp.pin()
		//sync_runtime_procPin()
		sync_runtime_procUnpin()
	}
}

func Benchmark_BuiltinAlloc(b *testing.B) {
	var mb []byte
	for i := 0; i < b.N; i++ {
		mb = make([]byte, 1024)
		if len(mb) != 0 {
		}
	}
}
