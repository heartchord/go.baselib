package goblazer

import "testing"
import "sort"

func Test_BinarySearchIntsGE01(t *testing.T) {
	var idx int
	var bytes = []int{ // 奇数个
		1 * 8,     // 1  byte
		2 * 8,     // 2  byte
		3 * 8,     // 3  byte
		4 * 8,     // 4  byte
		5 * 8,     // 5  byte
		6 * 8,     // 6  byte
		7 * 8,     // 7  byte
		8 * 8,     // 8  byte
		16 * 8,    // 16 byte
		32 * 8,    // 32 byte
		64 * 8,    // 64 byte
		96 * 8,    // 96 byte
		1 * 1024,  // 1  KB
		2 * 1024,  // 2  KB
		3 * 1024,  // 3  KB
		4 * 1024,  // 4  KB
		5 * 1024,  // 5  KB
		6 * 1024,  // 6  KB
		7 * 1024,  // 7  KB
		8 * 1024,  // 8  KB
		16 * 1024, // 16 KB
		32 * 1024, // 32 KB
		64 * 1024, // 64 KB
	}
	idx = BinarySearchIntsGE(bytes, 7)
	if idx != 0 {
		t.Error("Test_BinarySearchIntsGE01-01 failed")
		return
	}

	idx = BinarySearchIntsGE(bytes, 8)
	if idx != 0 {
		t.Error("Test_BinarySearchIntsGE01-02 failed")
		return
	}

	idx = BinarySearchIntsGE(bytes, 96*8)
	if idx != 11 {
		t.Error("Test_BinarySearchIntsGE01-03 failed")
		return
	}

	idx = BinarySearchIntsGE(bytes, 96*7)
	if idx != 11 {
		t.Error("Test_BinarySearchIntsGE01-04 failed")
		return
	}

	idx = BinarySearchIntsGE(bytes, 64*1024)
	if idx != 22 {
		t.Error("Test_BinarySearchIntsGE01-05 failed")
		return
	}

	idx = BinarySearchIntsGE(bytes, 64*1024-1)
	if idx != 22 {
		t.Error("Test_BinarySearchIntsGE01-06 failed")
		return
	}

	idx = BinarySearchIntsGE(bytes, 64*1024+1)
	if idx != -1 {
		t.Error("Test_BinarySearchIntsGE01-07 failed")
		return
	}

	t.Log("Test_BinarySearchIntsGE01 succeeded")
}

func Test_BinarySearchIntsGE02(t *testing.T) {
	var idx int
	var bytes = []int{ // 偶数个
		1 * 8,     // 1  byte
		2 * 8,     // 2  byte
		3 * 8,     // 3  byte
		4 * 8,     // 4  byte
		5 * 8,     // 5  byte
		6 * 8,     // 6  byte
		7 * 8,     // 7  byte
		8 * 8,     // 8  byte
		16 * 8,    // 16 byte
		32 * 8,    // 32 byte
		64 * 8,    // 64 byte
		96 * 8,    // 96 byte
		1 * 1024,  // 1  KB
		2 * 1024,  // 2  KB
		3 * 1024,  // 3  KB
		4 * 1024,  // 4  KB
		5 * 1024,  // 5  KB
		6 * 1024,  // 6  KB
		7 * 1024,  // 7  KB
		8 * 1024,  // 8  KB
		16 * 1024, // 16 KB
		32 * 1024, // 32 KB
	}

	idx = BinarySearchIntsGE(bytes, 7)
	if idx != 0 {
		t.Error("Test_BinarySearchIntsGE02-01 failed")
		return
	}

	idx = BinarySearchIntsGE(bytes, 8)
	if idx != 0 {
		t.Error("Test_BinarySearchIntsGE02-02 failed")
		return
	}

	idx = BinarySearchIntsGE(bytes, 96*8)
	if idx != 11 {
		t.Error("Test_BinarySearchIntsGE02-03 failed")
		return
	}

	idx = BinarySearchIntsGE(bytes, 96*7)
	if idx != 11 {
		t.Error("Test_BinarySearchIntsGE02-04 failed")
		return
	}

	idx = BinarySearchIntsGE(bytes, 32*1024)
	if idx != 21 {
		t.Error("Test_BinarySearchIntsGE02-05 failed")
		return
	}

	idx = BinarySearchIntsGE(bytes, 32*1024-1)
	if idx != 21 {
		t.Error("Test_BinarySearchIntsGE02-06 failed")
		return
	}

	idx = BinarySearchIntsGE(bytes, 64*1024)
	if idx != -1 {
		t.Error("Test_BinarySearchIntsGE02-07 failed")
		return
	}

	t.Log("Test_BinarySearchIntsGE02 succeeded")
}

func Test_BinarySearchIntsGE03(t *testing.T) {
	var idx int
	var bytes = []int{}

	idx = BinarySearchIntsGE(bytes, 1)
	if idx != -1 {
		t.Error("Test_BinarySearchIntsGE03 failed")
		return
	}

	t.Log("Test_BinarySearchIntsGE03 succeeded")
}

func Benchmark_BinarySearchIntsGE(b *testing.B) {
	b.StopTimer()
	var bytes = []int{ // 奇数个
		1 * 8,     // 1  byte
		2 * 8,     // 2  byte
		3 * 8,     // 3  byte
		4 * 8,     // 4  byte
		5 * 8,     // 5  byte
		6 * 8,     // 6  byte
		7 * 8,     // 7  byte
		8 * 8,     // 8  byte
		16 * 8,    // 16 byte
		32 * 8,    // 32 byte
		64 * 8,    // 64 byte
		96 * 8,    // 96 byte
		1 * 1024,  // 1  KB
		2 * 1024,  // 2  KB
		3 * 1024,  // 3  KB
		4 * 1024,  // 4  KB
		5 * 1024,  // 5  KB
		6 * 1024,  // 6  KB
		7 * 1024,  // 7  KB
		8 * 1024,  // 8  KB
		16 * 1024, // 16 KB
		32 * 1024, // 32 KB
		64 * 1024, // 64 KB
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		BinarySearchIntsGE(bytes, 1024)
	}
}

func Benchmark_BuiltinSearchInts(b *testing.B) {
	b.StopTimer()
	var bytes = []int{ // 奇数个
		1 * 8,     // 1  byte
		2 * 8,     // 2  byte
		3 * 8,     // 3  byte
		4 * 8,     // 4  byte
		5 * 8,     // 5  byte
		6 * 8,     // 6  byte
		7 * 8,     // 7  byte
		8 * 8,     // 8  byte
		16 * 8,    // 16 byte
		32 * 8,    // 32 byte
		64 * 8,    // 64 byte
		96 * 8,    // 96 byte
		1 * 1024,  // 1  KB
		2 * 1024,  // 2  KB
		3 * 1024,  // 3  KB
		4 * 1024,  // 4  KB
		5 * 1024,  // 5  KB
		6 * 1024,  // 6  KB
		7 * 1024,  // 7  KB
		8 * 1024,  // 8  KB
		16 * 1024, // 16 KB
		32 * 1024, // 32 KB
		64 * 1024, // 64 KB
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sort.SearchInts(bytes, 1024)
	}
}
