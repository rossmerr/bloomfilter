package bloomfilter

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/rossmerr/bitvector"
)

type HashFunction func([]byte) uint

type Filter[T Hash] struct {
	length uint
	vector bitvector.BitVector
	hash   HashFunction
}

type FilterOption[T Hash] func(*Filter[T])

func WithHash[T Hash](hash HashFunction) FilterOption[T] {
	return func(s *Filter[T]) {
		s.hash = hash
	}

}

// New Bloom filter
// n The number of elements in the filter
// p The probabilty of false postives
// m The number of elements in the BitVector
// k The number of hash functions to use
func NewFilter[T Hash](n uint, p float64, m, k uint, opts ...FilterOption[T]) *Filter[T] {
	if n < 1 {
		panic("capacity must be > 0")
	}
	if p >= 1 || p <= 0 {
		panic(fmt.Sprintf("probabity must be between 0 and 1, exclusive. Was %v", p))
	}
	if m < 1 {
		panic(fmt.Sprintf("The provided capacity and probabity values would result in an array of length > math.MaxInt32. Please reduce either of these values. Capacity: %v, probabity: %v", n, p))
	}

	var hasher = sha1.New()
	hash := func(data []byte) uint {
		hasher.Write([]byte(data))
		hash := hasher.Sum(nil)
		hasher.Reset()
		return uint(binary.LittleEndian.Uint64(hash))
	}

	filter := &Filter[T]{
		length: k,
		vector: *bitvector.NewBitVector(int(m)),
		hash:   hash,
	}

	for _, opt := range opts {
		opt(filter)
	}

	return filter
}

// New Bloom filter using the optimal size based on the capacity and probabity
func NewFilterOptimalWithProbabity[T Hash](n uint, p float64, opts ...FilterOption[T]) *Filter[T] {
	m := bestM(n, p)
	return NewFilter(n, p, m, bestK(n, m), opts...)
}

// New Bloom filter using the optimal size based on the capacity
func NewFilterOptimal[T Hash](n uint, opts ...FilterOption[T]) *Filter[T] {
	return NewFilterOptimalWithProbabity(n, probabity(n), opts...)
}

func probabity(n uint) float64 {
	c := 1.0 / float64(n)
	if c != 0 {
		return float64(c)
	}

	return math.Pow(0.6185, float64(math.MaxInt32/n))
}

// the number of bits
func bestM(n uint, p float64) uint {
	return uint(math.Round(-float64(n) * math.Log(p) / (math.Pow(math.Log(2), 2))))
}

// the number of hash functions
func bestK(n, m uint) uint {
	return uint(math.Round(float64(m) / float64(n) * math.Log(2)))
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
	firstHash, secondaryHash := s.splitHash(item)

	for i := uint(0); i < s.length; i++ {
		hash := s.computeHash(firstHash, secondaryHash, i)
		s.vector.Set(hash, true)
	}
}

func (s *Filter[T]) Length(item T) int {
	return int(s.length)
}

func (s *Filter[T]) computeHash(hash, secondaryHash, i uint) int {
	resultingHash := math.Mod(float64(hash)+(float64(i)*float64(secondaryHash)), float64(s.vector.Length()))
	return int(math.Abs(resultingHash))
}

// Checks for the existance of the item in the filter for a given probability.
func (s *Filter[T]) Contains(item T) bool {
	firstHash, secondaryHash := s.splitHash(item)

	for i := uint(0); i < s.length; i++ {
		hash := s.computeHash(firstHash, secondaryHash, i)
		if !s.vector.Get(hash) {
			return false
		}
	}

	return true
}

func (s *Filter[T]) splitHash(item T) (uint, uint) {
	hash := item.Sum()
	arr := intToBytes(hash)
	middle := int(math.Ceil(float64(len(arr)) / 2))
	first := arr[:middle]
	second := arr[middle:]
	firstHash := s.hash(first)
	secondaryHash := s.hash(second)
	return firstHash, secondaryHash
}

func intToBytes(num uint) []byte {
	buff := new(bytes.Buffer)
	bigOrLittleEndian := binary.LittleEndian
	err := binary.Write(buff, bigOrLittleEndian, uint64(num))
	if err != nil {
		panic(err)
	}

	return buff.Bytes()
}
