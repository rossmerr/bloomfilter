package bloomfilter

import (
	"fmt"
	"math"

	"github.com/rossmerr/bitvector"
)

type HashFunction[T Hash] func(T) int

type Hash interface {
	Sum() int
}

type Filter[T Hash] struct {
	count  int
	vector bitvector.BitVector
	hash   HashFunction[T]
}

// New Bloom filter
// m The number of elements in the BitVector
// k The number of hash functions to use
func NewFilter[T Hash](capacity int, errorRate float64, hash HashFunction[T], m, k int) *Filter[T] {
	if capacity < 1 {
		panic("capacity must be > 0")
	}
	if errorRate >= 1 || errorRate <= 0 {
		panic(fmt.Sprintf("errorRate must be between 0 and 1, exclusive. Was %v", errorRate))
	}
	if m < 1 {
		panic(fmt.Sprintf("The provided capacity and errorRate values would result in an array of length > int.MaxValue. Please reduce either of these values. Capacity: %v, Error rate: %v", capacity, errorRate))
	}

	return &Filter[T]{
		count:  k,
		vector: *bitvector.NewBitVector(m),
	}
}

// New Bloom filter using the optimal size based on the capacity and error rate
func NewFilterOptimalWithErrorRate[T Hash](capacity int, errorRate float64, hash HashFunction[T]) *Filter[T] {
	return NewFilter(capacity, errorRate, hash, int(bestM(capacity, errorRate)), int(bestK(capacity, errorRate)))
}

// New Bloom filter using the optimal size based on the capacity
func NewFilterOptimal[T Hash](capacity int, hash HashFunction[T]) *Filter[T] {
	return NewFilterOptimalWithErrorRate(capacity, bestErrorRate(capacity), hash)
}

func bestErrorRate(capacity int) float64 {
	c := 1.0 / capacity
	if c != 0 {
		return float64(c)
	}

	// default
	// http://www.cs.princeton.edu/courses/archive/spring02/cs493/lec7.pdf
	return math.Pow(0.6185, float64(math.MaxInt32/capacity))
}
func bestK(capacity int, errorRate float64) float64 {
	return math.Round(math.Log(2.0) * bestM(capacity, errorRate) / float64(capacity))
}

func bestM(capacity int, errorRate float64) float64 {
	//t := (1.0 / math.Pow(2, math.Log(2.0)))
	return math.Ceil(float64(capacity) * math.Log(errorRate))
}

// The number of true bits.
func (s *Filter[T]) TrueBits() int {
	output := 0
	iterator := s.vector.Enumerate()
	for iterator.HasNext() {
		v, _ := iterator.Next()
		if v {
			output++
		}
	}
	return output

}

// The ratio of false to true bits in the BitVector.
func (s *Filter[T]) Truthiness() float64 {
	return float64(s.TrueBits()) / float64(s.vector.Length())
}

// Adds a new item to the filter. It cannot be removed.
func (s *Filter[T]) Add(item T) {
	hash := item.Sum()
	secondaryHash := s.hash(item)
	for i := 0; i < s.count; i++ {
		hash := s.computeHash(hash, secondaryHash, i)
		s.vector.Set(hash, true)
	}
}

func (s *Filter[T]) computeHash(hash, secondaryHash, i int) int {
	resultingHash := math.Mod(float64(hash)+(float64(i)*float64(secondaryHash)), float64(s.vector.Length()))
	return int(math.Abs(resultingHash))
}

// Checks for the existance of the item in the filter for a given probability.
func (s *Filter[T]) Contains(item T) bool {
	hash := item.Sum()
	secondaryHash := s.hash(item)

	for i := 0; i < s.count; i++ {
		hash := s.computeHash(hash, secondaryHash, i)
		if !s.vector.Get(hash) {
			return false
		}
	}

	return true
}
