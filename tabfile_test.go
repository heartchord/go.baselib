package goblazer

import "testing"

func Benchmark_TabFileLoad_UTF8(b *testing.B) {
	b.StopTimer()

	f := NewTabFile()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		f.Load("C:/npcs_utf8.txt")
	}
}

func Benchmark_TabFileLoad_GBK(b *testing.B) {
	b.StopTimer()

	f := NewTabFile()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		f.Load("C:/npcs_gbk.txt")
	}
}

func Benchmark_TabFileSave_UTF8(b *testing.B) {
	b.StopTimer()

	f := NewTabFile()
	f.Load("C:/npcs_utf8.txt")

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		f.Save("C:/npcs_utf8_bak.txt")
	}
}

func Benchmark_TabFileSave_GBK(b *testing.B) {
	b.StopTimer()

	f := NewTabFile()
	f.Load("C:/npcs_gbk.txt")

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		f.Save("C:/npcs_gbk_bak.txt")
	}
}
