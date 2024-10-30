# go-bloom

A very simple fast bloom filter with Add, Test, Save, and Load methods.

## Create and test
```golang
  filter := bloom.Filter{make([]byte, 100)}
  filter.Add([]byte("hello"))
  hit := filter.TestString("hello")
```

To save or load data, use the filter.Data slice.

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
