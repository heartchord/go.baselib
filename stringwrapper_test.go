package goblazer

import (
	"strings"
	"testing"
)

func Test_BytesToStringByShare(t *testing.T) {
	var bytes = []byte("Whatever is worth doing is worth doing well.")

	s1 := string(bytes)
	s2 := BytesToStringByShare(bytes)

	if strings.Compare(s1, s2) == 0 {
		t.Log("Test_BytesToStringByShare succeeded")
	} else {
		t.Error("Test_BytesToStringByShare failed")
	}
}

func Test_BytesToStringByTrans(t *testing.T) {
	var bytes = []byte("Whatever is worth doing is worth doing well.")

	s1 := string(bytes)
	s2 := BytesToStringByTrans(&bytes)

	if strings.Compare(s1, s2) == 0 {
		t.Log("Test_BytesToStringByTrans succeeded")
	} else {
		t.Error("Test_BytesToStringByTrans failed")
	}
}

func Benchmark_BytesToStringByShare(b *testing.B) {
	var bytes = []byte("Whatever is worth doing is worth doing well.")

	for i := 0; i < b.N; i++ {
		_ = BytesToStringByShare(bytes)
	}
}

func Benchmark_BytesToStringByTrans(b *testing.B) {
	var bytes = []byte("Whatever is worth doing is worth doing well.")

	for i := 0; i < b.N; i++ {
		_ = BytesToStringByTrans(&bytes)
	}
}

func Benchmark_BytesToStringByCopys(b *testing.B) {
	var bytes = []byte("Whatever is worth doing is worth doing well.")

	for i := 0; i < b.N; i++ {
		_ = BytesToStringByCopys(bytes)
	}
}

func Test_StringToBytesByShare(t *testing.T) {
	var s = "Whatever is worth doing is worth doing well."
	var r = true

	b1 := []byte(s)
	b2 := StringToBytesByShare(s)

	if len(b1) != len(b2) {
		t.Error("Test_StringToBytesByShare failed")
	}

	for i, v1 := range b1 {
		if v1 != b2[i] {
			r = false
			break
		}
	}

	if r {
		t.Log("Test_StringToBytesByShare succeeded")
	} else {
		t.Error("Test_StringToBytesByShare failed")
	}
}

func Test_StringToBytesByTrans(t *testing.T) {
	var s = "Whatever is worth doing is worth doing well."
	var r = true

	b1 := []byte(s)
	b2 := StringToBytesByTrans(&s)

	if len(b1) != len(b2) {
		t.Error("Test_StringToBytesByTrans failed")
	}

	for i, v1 := range b1 {
		if v1 != b2[i] {
			r = false
			break
		}
	}

	if r {
		t.Log("Test_StringToBytesByTrans succeeded")
	} else {
		t.Error("Test_StringToBytesByTrans failed")
	}
}

func Benchmark_StringToBytesByShare(b *testing.B) {
	var s = "Whatever is worth doing is worth doing well."

	for i := 0; i < b.N; i++ {
		_ = StringToBytesByShare(s)
	}
}

func Benchmark_StringToBytesByTrans(b *testing.B) {
	var s = "Whatever is worth doing is worth doing well."

	for i := 0; i < b.N; i++ {
		_ = StringToBytesByTrans(&s)
	}
}

func Benchmark_StringToBytesByCopys(b *testing.B) {
	var s = "Whatever is worth doing is worth doing well."

	for i := 0; i < b.N; i++ {
		_ = StringToBytesByCopys(s)
	}
}
