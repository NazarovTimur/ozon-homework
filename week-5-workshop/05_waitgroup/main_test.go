package main

import (
	"testing"
)

// go test --bench=.

func Benchmark_Simple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iterate()
	}
}

func Benchmark_Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iterateInt()
	}
}

func Benchmark_Atomic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iterateAtomic()
	}
}

func Benchmark_Mutext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iterateMutex()
	}
}
