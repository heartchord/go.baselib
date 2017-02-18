package goblazer

import (
	"encoding/binary"
	"testing"
)

func Test_BytesToNumber_1(t *testing.T) {
	var n uint16

	b := []byte{0x7E, 0x7F}
	if ok := BytesToNumber(b, binary.LittleEndian, &n); !ok {
		t.Error("Test_BytesToNumber_1 failed")
	}

	r := uint16(32638)
	if n == r {
		t.Log("Test_BytesToNumber_1 succeeded")
	} else {
		t.Error("Test_BytesToNumber_1 failed")
	}
}

func Test_BytesToNumber_2(t *testing.T) {
	var n uint32

	b := []byte{0xFF, 0x7E, 0xFE, 0x7F}
	if ok := BytesToNumber(b, binary.LittleEndian, &n); !ok {
		t.Error("Test_BytesToNumber_2 failed")
	}

	r := uint32(2147385087)
	if n == r {
		t.Log("Test_BytesToNumber_2 succeeded")
	} else {
		t.Error("Test_BytesToNumber_2 failed")
	}
}

func Test_BytesToUint32_1(t *testing.T) {
	b := []byte{0xFF, 0x7E, 0xFE, 0x7F}
	n := BytesToUint32(b)

	r := uint32(2147385087)
	if n == r {
		t.Log("Test_BytesToUint32_1 succeeded")
	} else {
		t.Error("Test_BytesToUint32_1 failed")
	}
}

func Benchmark_BytesToNumber_1(b *testing.B) {
	var n uint16
	bs := []byte{0x7E, 0x7F}

	for i := 0; i < b.N; i++ {
		BytesToNumber(bs, binary.LittleEndian, &n)
	}
}

func Benchmark_BytesToNumber_2(b *testing.B) {
	var n uint32
	bs := []byte{0xFF, 0x7E, 0xFE, 0x7F}

	for i := 0; i < b.N; i++ {
		BytesToNumber(bs, binary.LittleEndian, &n)
	}
}

func Benchmark_BytesToUint32_1(b *testing.B) {
	bs := []byte{0xFF, 0x7E, 0xFE, 0x7F}

	for i := 0; i < b.N; i++ {
		BytesToUint32(bs)
	}
}
