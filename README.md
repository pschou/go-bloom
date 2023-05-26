# go-bloom

A very simple fast bloom filter with Add, Test, Save, and Load methods.

## Create and test
```golang
  filter := New(100)
  filter.AddString("hello")
  fmt.Println("test", filter.TestString("hello"))
  // Output:
  // test true
```

## Save and loading
```golang
  filter := New(100)
  filter.Add([]byte("hello"))
  fh, _ := os.Create("bloom.flt")
  filter.Save(fh)
  fh.Close()

  my, _ := os.Open("bloom.flt")
  filt, _ := Load(my)
  my.Close()
  fmt.Println("test", filt.TestString("hello"))
  // Output:
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
