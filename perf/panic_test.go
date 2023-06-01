package perf

import (
	"fmt"
	"testing"
)

func BenchmarkPanic(b *testing.B) {
	pf := func() {
		panic(1)
		fmt.Println()
	}
	f := func() {
		defer func() {
			_ = recover()
		}()

		pf()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f()
	}
}

func BenchmarkEarlyReturn(b *testing.B) {
	pf := func() {
		return
		fmt.Println()
	}
	f := func() {
		defer func() {
			_ = recover()
		}()

		pf()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f()
	}
}
