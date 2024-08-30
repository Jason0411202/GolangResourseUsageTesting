package main

import (
	"testing"
)

func BenchmarkCallByValue(b *testing.B) {
	for n := 0; n < b.N; n++ {
		CallByValueTest()
	}
}

func BenchmarkCallByReference(b *testing.B) {
	for n := 0; n < b.N; n++ {
		CallByReferenceTest()
	}
}
