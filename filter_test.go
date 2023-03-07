package bloomfilter_test

import (
	"testing"

	"github.com/rossmerr/bloomfilter"
)

type test struct {
}

func (s *test) Sum() int {
	return 1
}

func TestFilter_Add_Contains(t *testing.T) {
	tests := []struct {
		name   string
		values *test
		hash   bloomfilter.HashFunction[*test]
		want   bool
	}{
		{
			name:   "Add",
			values: &test{},
			hash: func(t *test) int {
				return 2
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := bloomfilter.NewFilterOptimalWithProbabity(216553, 0.01, tt.hash)

			filter.Add(tt.values)
			got := filter.Contains(tt.values)
			if got != tt.want {
				t.Errorf("Filter.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_Contains(t *testing.T) {
	tests := []struct {
		name   string
		values *test
		hash   bloomfilter.HashFunction[*test]
		want   bool
	}{
		{
			name:   "Does not contains",
			values: &test{},
			hash: func(t *test) int {
				return 2
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := bloomfilter.NewFilterOptimalWithProbabity(216553, 0.01, tt.hash)

			got := filter.Contains(tt.values)
			if got != tt.want {
				t.Errorf("Filter.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
