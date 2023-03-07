package bloomfilter

import (
	"fmt"
	"math"

	"github.com/rossmerr/bitvector"
)

type HashFunction[T Hash] func(T) int

type Filter[T Hash] struct {
	count  int
	vector bitvector.BitVector
	hash   HashFunction[T]
}

// New Bloom filter
// n The number of elements in the filter
// p The probabilty of false postives
// m The number of elements in the BitVector
// k The number of hash functions to use
func NewFilter[T Hash](n int, p float64, hash HashFunction[T], m, k int) *Filter[T] {
	if n < 1 {
		panic("capacity must be > 0")
	}
	if p >= 1 || p <= 0 {
		panic(fmt.Sprintf("probabity must be between 0 and 1, exclusive. Was %v", p))
	}
	if m < 1 {
		panic(fmt.Sprintf("The provided capacity and probabity values would result in an array of length > int.MaxValue. Please reduce either of these values. Capacity: %v, probabity: %v", n, p))
	}

	return &Filter[T]{
		count:  k,
		vector: *bitvector.NewBitVector(m),
		hash:   hash,
	}
}

// New Bloom filter using the optimal size based on the capacity and probabity
func NewFilterOptimalWithProbabity[T Hash](n int, p float64, hash HashFunction[T]) *Filter[T] {
	m := bestM(n, p)
	return NewFilter(n, p, hash, m, bestK(n, m))
}

// New Bloom filter using the optimal size based on the capacity
func NewFilterOptimal[T Hash](n int, hash HashFunction[T]) *Filter[T] {
	return NewFilterOptimalWithProbabity(n, probabity(n), hash)
}

func probabity(n int) float64 {
	c := 1.0 / float64(n)
	if c != 0 {
		return float64(c)
	}

	return math.Pow(0.6185, float64(math.MaxInt32/n))
}

// the number of bits
func bestM(n int, p float64) int {
	return int(math.Round(-float64(n) * math.Log(p) / (math.Pow(math.Log(2), 2))))
}

// the number of hash functions
func bestK(n, m int) int {
	return int(math.Round(float64(m) / float64(n) * math.Log(2)))
}

// The number of true bits.
func (s *Filter[T]) TrueBits() int {
	return s.vector.TrueBits()
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
