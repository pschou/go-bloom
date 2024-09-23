package bloom_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/pschou/go-bloom"
)

func ExampleNew() {
	filter := bloom.New(100)
	filter.Add([]byte("hello"))
	hit := filter.Test([]byte("hello"))
	fmt.Println("test", hit)
	// Output:
	// test true
}

func ExampleAddString() {
	filter := bloom.New(100)
	filter.AddString("hello")
	hit := filter.TestString("hello")
	fmt.Println("test", hit)
	// Output:
	// test true
}

func ExampleAdd() {
	filter := bloom.New(100)
	filter.Add([]byte("hello"))
	hit := filter.TestString("hello")
	fmt.Println("test", hit)
	// Output:
	// test true
}

func ExampleSaveAndLoad() {
	filter := bloom.New(100)
	filter.Add([]byte("hello"))
	fh, _ := os.Create("bloom.flt")
	filter.Save(fh)
	fh.Close()

	filt, err := bloom.LoadFile("bloom.flt", 1)
	if err != nil {
		panic(err)
	}

	// Fold the filter in half
	filt.Fold(2)

	// Verify that hello still matches
	hit := filt.TestString("hello")
	fmt.Println("test", hit)
	// Output:
	// test true
}

func ExampleSaveAndLoad5() {
	filter := bloom.New(100)
	filter.Add([]byte("hello"))
	fh, _ := os.Create("bloom5.flt")
	filter.Save(fh)
	fh.Close()

	filt, err := bloom.LoadFile("bloom5.flt", 5)
	if err != nil {
		panic(err)
	}

	// Fold the filter in half
	filt.Fold(2)

	// Verify that hello still matches
	hit := filt.TestString("hello")
	fmt.Println("test", hit)
	miss := filt.TestString("hello2")
	fmt.Println("test", miss)
	// Output:
	// test true
	// test false
}

func BenchmarkAdd(b *testing.B) {
	dat := []byte("helloworld")
	filter := bloom.New(1 << 24)
	for n := 0; n < b.N; n++ {
		filter.Add(dat)
	}
}

func BenchmarkAddString(b *testing.B) {
	dat := "helloworld"
	filter := bloom.New(1 << 24)
	for n := 0; n < b.N; n++ {
		filter.AddString(dat)
	}
}

func BenchmarkTest(b *testing.B) {
	dat := []byte("helloworld")
	filter := bloom.New(1 << 24)
	for n := 0; n < b.N; n++ {
		filter.Test(dat)
	}
}

func BenchmarkTestString(b *testing.B) {
	dat := "helloworld"
	filter := bloom.New(1 << 24)
	for n := 0; n < b.N; n++ {
		filter.TestString(dat)
	}
}
