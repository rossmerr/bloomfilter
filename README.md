# Bloom Filter

## Usage

define your struct implemting the `Hash` interface defined in this module.

```go
type Test struct {

}

func (s *Test) Sum() int {
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
