package bloom

import (
	"fmt"
	"os"
	"testing"
)

func ExampleNew() {
	filter := New(100)
	fmt.Println("before", filter.dat)
	filter.AddString("hello")
	fmt.Println("after", filter.dat)
	_, hit := filter.TestString("hello")
	fmt.Println("test", hit)
	// Output:
	// before [0 0 0 0 0 0 0 0 0 0 0 0]
	// after [0 0 0 32 0 0 0 0 0 0 0 0]
	// test true
}

func ExampleAddString() {
	filter := New(100)
	fmt.Println("before", filter.dat)
	filter.AddString("hello")
	fmt.Println("after", filter.dat)
	_, hit := filter.TestString("hello")
	fmt.Println("test", hit)
	// Output:
	// before [0 0 0 0 0 0 0 0 0 0 0 0]
	// after [0 0 0 32 0 0 0 0 0 0 0 0]
	// test true
}

func ExampleAdd() {
	filter := New(100)
	fmt.Println("before", filter.dat)
	filter.Add([]byte("hello"))
	fmt.Println("after", filter.dat)
	_, hit := filter.TestString("hello")
	fmt.Println("test", hit)
	// Output:
	// before [0 0 0 0 0 0 0 0 0 0 0 0]
	// after [0 0 0 32 0 0 0 0 0 0 0 0]
	// test true
}

func ExampleSaveAndLoad() {
	filter := New(100)
	filter.Add([]byte("hello"))
	fmt.Println("before", filter.dat)
	fh, _ := os.Create("bloom.flt")
	filter.Save(fh)
	fh.Close()

	my, _ := os.Open("bloom.flt")
	filt, _ := Load(my)
	my.Close()
	fmt.Println("after", filt.dat)
	_, hit := filter.TestString("hello")
	fmt.Println("test", hit)
	// Output:
	// before [0 0 0 32 0 0 0 0 0 0 0 0]
	// after [0 0 0 32 0 0 0 0 0 0 0 0]
	// test true
}

func BenchmarkAdd(b *testing.B) {
	dat := []byte("helloworld")
	filter := New(1 << 32)
	for n := 0; n < b.N; n++ {
		filter.Add(dat)
	}
}

func BenchmarkAddString(b *testing.B) {
	dat := "helloworld"
	filter := New(1 << 32)
	for n := 0; n < b.N; n++ {
		filter.AddString(dat)
	}
}

func BenchmarkTest(b *testing.B) {
	dat := []byte("helloworld")
	filter := New(1 << 32)
	for n := 0; n < b.N; n++ {
		filter.Test(dat)
	}
}

func BenchmarkTestString(b *testing.B) {
	dat := "helloworld"
	filter := New(1 << 32)
	for n := 0; n < b.N; n++ {
		filter.TestString(dat)
	}
}
