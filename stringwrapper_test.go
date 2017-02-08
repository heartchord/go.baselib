package goblazer

import (
	"fmt"
	"strings"
	"testing"
)

// Test - BytesToString
func Test_BytesToString(t *testing.T) {
	var bytes = []byte("Whatever is worth doing is worth doing well.")

	s1 := string(bytes)
	s2 := BytesToString(bytes)

	if strings.Compare(s1, s2) == 0 {
		t.Log("Test_BytesToString succeeded")
	} else {
		t.Error("Test_BytesToString failed")
	}
}

// Test - StringToBytes
func Test_StringToBytes(t *testing.T) {
	var s = "Whatever is worth doing is worth doing well."
	var r = true

	b1 := []byte(s)
	b2 := StringToBytes(s)

	fmt.Println(cap(b1))
	fmt.Println(cap(b2))

	if len(b1) != len(b2) {
		t.Error("Test_StringToBytes failed")
		return
	}

	for i, v1 := range b1 {
		if v1 != b2[i] {
			r = false
			break
		}
	}

	if r {
		t.Log("Test_StringToBytes succeeded")
	} else {
		t.Error("Test_StringToBytes failed")
	}
}

// Benchmark - BytesToString
func Benchmark_BytesToString(b *testing.B) {
	var bytes = []byte("Whatever is worth doing is worth doing well.")

	for i := 0; i < b.N; i++ {
		_ = BytesToString(bytes)
	}
}

// Benchmark - BytesToStringDefault
func Benchmark_BytesToStringDefault(b *testing.B) {
	var bytes = []byte("Whatever is worth doing is worth doing well.")

	for i := 0; i < b.N; i++ {
		_ = string(bytes)
	}
}

// Benchmark - StringToBytes
func Benchmark_StringToBytes(b *testing.B) {
	var s = "Whatever is worth doing is worth doing well."

	for i := 0; i < b.N; i++ {
		_ = StringToBytes(s)
	}
}

func Benchmark_StringToBytesDefault(b *testing.B) {
	var s = "Whatever is worth doing is worth doing well."

	for i := 0; i < b.N; i++ {
		_ = []byte(s)
	}
}
