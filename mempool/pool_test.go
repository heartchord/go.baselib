// +build !race

package mempool

import (
	"runtime"
	"sync"
	"testing"
)

func testNewFunc() interface{} {
	return 1234567890
}

func Test_Pool(t *testing.T) {
	p := NewPool(2, nil)
	p.EnableStats()

	x := p.Get()
	if x != nil {
		t.Fatalf("expected nil, got %v", x)
	}

	p.Put(nil)
	p.Put(100)
	p.Put(200)
	p.Put(300)

	v0 := p.Get()
	if i0 := v0.(int); i0 != 100 {
		t.Fatalf("expected 100, got %v", i0)
	}

	v1 := p.Get()
	if i1 := v1.(int); i1 != 200 {
		t.Fatalf("expected 200, got %v", i1)
	}

	v2 := p.Get()
	if v2 != nil {
		t.Fatalf("expected nil, got %v", v2)
	}

	p.GetStats().PrintAllInfo()
}

func Test_PoolNew(t *testing.T) {
	p := NewPool(2, testNewFunc)
	p.EnableStats()

	x := p.Get()
	if i := x.(int); i != 1234567890 {
		t.Fatalf("expected 1234567890, got %v", i)
	}
	p.GetStats().PrintAllInfo()
}

func Test_PoolReset(t *testing.T) {
	p := NewPool(2, nil)
	p.EnableStats()

	p.Put(20)
	p.Reset()

	v0 := p.Get()
	if v0 != nil {
		t.Fatalf("expected nil, got %v", v0)
	}
	p.GetStats().PrintAllInfo()
}

func Test_PoolStress(t *testing.T) {
	const P = 10
	N := int(1e6)
	if testing.Short() {
		N /= 100
	}

	p := NewPool(100, nil)
	p.EnableStats()
	done := make(chan bool)

	for i := 0; i < P; i++ {
		go func() {
			var v interface{}
			for j := 0; j < N; j++ {
				if v == nil {
					v = 0
				}

				p.Put(v)
				v = p.Get()
				if v != nil && v.(int) != 0 {
					t.Fatalf("expected 0, got %v", v)
				}
			}
			done <- true
		}()
	}

	for i := 0; i < P; i++ {
		<-done
	}
	p.GetStats().PrintAllInfo()
}

func Test_PoolPChanged(t *testing.T) {
	const P = 10
	N := int(1e4)
	numCPU := runtime.GOMAXPROCS(1)
	p := NewPool(100, nil)
	p.EnableStats()
	runtime.GOMAXPROCS(numCPU)

	done := make(chan bool)
	for i := 0; i < P; i++ {
		go func() {
			var v interface{}
			for j := 0; j < N; j++ {
				if v == nil {
					v = 0
				}
				p.Put(v)
				v = p.Get()
				if v != nil && v.(int) != 0 {
					t.Fatalf("expected 0, got %v", v)
				}
			}
			done <- true
		}()
	}

	for i := 0; i < P; i++ {
		<-done
	}
	p.GetStats().PrintAllInfo()
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

func Benchmark_BuiltinPool(b *testing.B) {
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
