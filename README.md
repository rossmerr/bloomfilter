# Bloom Filter

[![Go](https://github.com/rossmerr/bloomfilter/actions/workflows/go.yml/badge.svg)](https://github.com/rossmerr/bloomfilter/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rossmerr/bloomfilter)](https://goreportcard.com/report/github.com/rossmerr/bloomfilter)
[![Read the Docs](https://pkg.go.dev/badge/golang.org/x/pkgsite)](https://pkg.go.dev/github.com/rossmerr/bloomfilter)

## Usage

define your struct implementing the [Hash](hash.go) interface defined in this module.

```go
type Test struct {

}

func (s *Test) Sum() uint {
  // your hash function...
}
```

```go
obj := &Test{}

filter := bloomfilter.NewFilterOptimal[Test](2000000)
filter.Add(obj)

match := filter.Contains(obj)
fmt.Println(match) // true
```
