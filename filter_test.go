package bloomfilter_test

import (
	"testing"

	"github.com/rossmerr/bloomfilter"
)

type test struct {
}

func (s *test) Sum() uint {
	return 10
}

func TestFilter_Add_Contains(t *testing.T) {
	tests := []struct {
		name   string
		values *test
		want   bool
	}{
		{
			name:   "Add",
			values: &test{},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := bloomfilter.NewFilterOptimalWithProbabity[*test](216553, 0.01)

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
		want   bool
	}{
		{
			name:   "Does not contains",
			values: &test{},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := bloomfilter.NewFilterOptimalWithProbabity[*test](216553, 0.01)

			got := filter.Contains(tt.values)
			if got != tt.want {
				t.Errorf("Filter.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
