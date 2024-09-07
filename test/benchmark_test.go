package test

import (
	"bytes"
	"fmt"
	"io"
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

// 3 functions writing to buffer
func writeToBufferByVal(buf bytes.Buffer, data []byte) {
	buf.Write(data)
}

func writeToBufferByPtr(buff *bytes.Buffer, data []byte) {
	buff.Write(data)
}

func writeToWriter(writer io.Writer, data []byte) {
	writer.Write(data)
}

func BenchmarkWriteToBufferByVal(b *testing.B) {
	buf := bytes.Buffer{}
	data := []byte("data")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writeToBufferByVal(buf, data)
	}
}

func BenchmarkWriteToBufferByPtr(b *testing.B) {
	buf := new(bytes.Buffer)
	data := []byte("data")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writeToBufferByPtr(buf, data)
	}
}

func BenchmarkWriteToBufferByPtrWithReset(b *testing.B) {
	buf := new(bytes.Buffer)
	data := []byte("data")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		writeToBufferByPtr(buf, data)
	}
}

func BenchmarkWriteToWriter(b *testing.B) {
	buf := new(bytes.Buffer)
	data := []byte("data")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writeToWriter(buf, data)
	}
}

func BenchmarkWriteToWriterWithReset(b *testing.B) {
	buf := new(bytes.Buffer)
	data := []byte("data")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		writeToWriter(buf, data)
	}
}
