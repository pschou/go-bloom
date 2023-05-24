# go-bloom

A very simple fast bloom filter with Add, Test, Save, and Load methods.

## Create and test
```golang
  filter := New(100)
  fmt.Println("before", filter.dat)
  filter.AddString("hello")
  fmt.Println("after", filter.dat)
  fmt.Println("test", filter.TestString("hello"))
  // Output:
  // before [0 0 0 0 0 0 0 0 0 0 0 0]
  // after [0 0 0 32 0 0 0 0 0 0 0 0]
  // test true
```

## Save and loading
```golang
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
  fmt.Println("test", filt.TestString("hello"))
  // Output:
  // before [0 0 0 32 0 0 0 0 0 0 0 0]
  // after [0 0 0 32 0 0 0 0 0 0 0 0]
  // test true
```

## Benchmarks
```
$ go test --bench=.
goos: linux
goarch: amd64
cpu: Intel(R) Xeon(R) CPU           X5650  @ 2.67GHz
BenchmarkAdd-12                 40322137                26.81 ns/op
BenchmarkAddString-12           32753798                38.91 ns/op
BenchmarkTest-12                44796183                24.68 ns/op
BenchmarkTestString-12          40576682                27.15 ns/op
PASS
```
