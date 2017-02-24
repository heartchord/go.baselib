package mempool

import "testing"

func Test_BytesPool(t *testing.T) {
	b0 := GetFromDefaultBytesPool(100)
	if b0 == nil {
		t.Fatal("b0 should not be nil")
	}
	copy(b0.Buf(), "hello")
	if len(b0.Buf()) != 100 {
		t.Fatalf("expected 100, got %v", len(b0.Buf()))
	}
	if cap(b0.Buf()) != 128 {
		t.Fatalf("expected 128, got %v", cap(b0.Buf()))
	}
	b0.AddRef()
	b0.DecRef()
	b0.DecRef()

	b1 := GetFromDefaultBytesPool(32)
	if b1 == nil {
		t.Fatal("b1 should not be nil")
	}
	if len(b1.Buf()) != 32 {
		t.Fatalf("expected 32, got %v", len(b1.Buf()))
	}
	if cap(b1.Buf()) != 32 {
		t.Fatalf("expected 32, got %v", cap(b1.Buf()))
	}
	b1.DecRef()

	// 超过最大的class的size
	b2 := GetFromDefaultBytesPool(16 * 1024 * 1024)
	if b2 == nil {
		t.Fatal("b2 should not be nil")
	}
	if len(b2.Buf()) != 16*1024*1024 {
		t.Fatalf("expected 16*1024*1024, got %v", len(b2.Buf()))
	}
	if cap(b2.Buf()) != 16*1024*1024 {
		t.Fatalf("expected 16*1024*1024, got %v", cap(b2.Buf()))
	}
	b2.DecRef()

	// 请求0字节
	b3 := GetFromDefaultBytesPool(0)
	if b3 == nil {
		t.Fatal("b3 should not be nil")
	}
	if len(b3.Buf()) != 0 {
		t.Fatalf("expected 0, got %v", len(b3.Buf()))
	}
	if cap(b3.Buf()) != 16 {
		t.Fatalf("expected 16, got %v", cap(b3.Buf()))
	}
	b3.DecRef()
}

func Test_InvalidClasses(t *testing.T) {
	if checkAllClasses([]int{}) {
		t.Fatal("checkClasses should return false")
	}

	classes0 := [...]int{0, 1, 2, 3}
	if checkAllClasses(classes0[:]) {
		t.Fatal("checkClasses should return false")
	}

	classes1 := [...]int{1, 2, 2, 4}
	if checkAllClasses(classes1[:]) {
		t.Fatal("checkClasses should return false")
	}

	classes2 := [...]int{1, 2, 8, 4}
	if checkAllClasses(classes2[:]) {
		t.Fatal("checkClasses should return false")
	}
}

func Test_Clear(t *testing.T) {
	ResetDefaultBytesPool()
	b0 := GetFromDefaultBytesPool(100)
	b1 := GetFromDefaultBytesPool(255)
	b1.DecRef()
	b0.DecRef()
	ResetDefaultBytesPool()
}

func Benchmark_BytesPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := GetFromDefaultBytesPool(1024)
			buf.DecRef()
		}
	})
}

func Benchmark_BytesPoolLargeSize(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := GetFromDefaultBytesPool(65536 + 1)
			buf.DecRef()
		}
	})
}

func Benchmark_BuiltinAlloc(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var dummy []byte
		for pb.Next() {
			buf := make([]byte, 1024)
			dummy = buf // 防止buf的分配被优化掉
		}
		if len(dummy) > 0 {
			dummy[0] = 0
		}
	})
}

func Benchmark_BuiltinAllocLargeSize(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var dummy []byte
		for pb.Next() {
			buf := make([]byte, 65536+1)
			dummy = buf // 防止buf的分配被优化掉
		}
		if len(dummy) > 0 {
			dummy[0] = 0
		}
	})
}
