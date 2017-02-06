package goblazer

import "testing"

func Benchmark_TabFileLoad(b *testing.B) {
	f := NewTabFile()

	for i := 0; i < b.N; i++ {
		f.Load("C:\\Users\\chris\\Desktop\\bakfiles\\testtab.txt", "utf8")
	}
}
