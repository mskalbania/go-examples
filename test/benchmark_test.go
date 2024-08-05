package test

import (
	"fmt"
	"strconv"
	"testing"
)

// use -test.benchtime 5s to increase the time of benchmark
func BenchmarkSprintf(b *testing.B) {
	number := 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%d", number)
	}
}

func BenchmarkFormatInt(b *testing.B) {
	number := int64(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strconv.FormatInt(number, 10)
	}
}

func BenchmarkItoa(b *testing.B) {
	number := 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strconv.Itoa(number)
	}
}
