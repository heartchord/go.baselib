package mempool

import "testing"

func Benchmark_BytesPool_1(b *testing.B) {
	b.StopTimer()
	bp := newLocalBytesPool()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		bb := bp.getPrivate(1024)
		bp.putPrivate(bb)
		//getIndex(1024)
	}
}

func Benchmark_BytesPool(b *testing.B) {
	b.StopTimer()
	bp := NewBytesPool()
	//for i := 0; i < 100; i++ {
	//	go func() {
	//		for {
	//			bb := bp.Get(1024)
	//			time.Sleep(20)
	//			bp.Put(bb)
	//		}
	//	}()
	//}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		bb := bp.Get(1024)
		bp.Put(bb)
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
