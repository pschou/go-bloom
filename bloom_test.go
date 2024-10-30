package bwdb_test

import (
	"fmt"
	"testing"

	bloom "github.com/pschou/go-bloom"
)

func ExampleAddString() {
	filter := bloom.Filter{make([]byte, 100)}
	filter.AddString("hello")
	hit := filter.TestString("hello")
	fmt.Println("test", hit)
	// Output:
	// test true
}

func ExampleAdd() {
	filter := bloom.Filter{make([]byte, 100)}
	filter.Add([]byte("hello"))
	hit := filter.TestString("hello")
	fmt.Println("test", hit)
	// Output:
	// test true
}

func ExampleFold() {
	filter := bloom.Filter{make([]byte, 100)}
	filter.Add([]byte("hello"))

	// Fold the filter in half
	fmt.Println("before folding size:", len(filter.Data))
	filter.Fold(5)
	fmt.Println("after folding size:", len(filter.Data))

	// Verify that hello still matches
	hit := filter.TestString("hello")
	fmt.Println("test", hit)
	// Output:
	// before folding size: 100
	// after folding size: 20
	// test true
}

func BenchmarkAdd(b *testing.B) {
	dat := []byte("helloworld")
	filter := bloom.Filter{make([]byte, 1<<24)}
	for n := 0; n < b.N; n++ {
		filter.Add(dat)
	}
}

func BenchmarkAddString(b *testing.B) {
	dat := "helloworld"
	filter := bloom.Filter{make([]byte, 1<<24)}
	for n := 0; n < b.N; n++ {
		filter.AddString(dat)
	}
}

func BenchmarkTest(b *testing.B) {
	dat := []byte("helloworld")
	filter := bloom.Filter{make([]byte, 1<<24)}
	for n := 0; n < b.N; n++ {
		filter.Test(dat)
	}
}

func BenchmarkTestString(b *testing.B) {
	dat := "helloworld"
	filter := bloom.Filter{make([]byte, 1<<24)}
	for n := 0; n < b.N; n++ {
		filter.TestString(dat)
	}
}
