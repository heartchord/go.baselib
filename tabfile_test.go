package goblazer

import "testing"

func Benchmark_TabFileLoad_UTF8(b *testing.B) {
	f := NewTabFile()

	for i := 0; i < b.N; i++ {
		f.Load("C:/npcs_utf8.txt", "utf8")
	}
}

func Benchmark_TabFileLoad_GBK(b *testing.B) {
	f := NewTabFile()

	for i := 0; i < b.N; i++ {
		f.Load("C:/npcs.txt", "gbk")
	}
}
