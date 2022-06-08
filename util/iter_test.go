package util

import (
	"testing"
)

func TestIter(t *testing.T) {
	xs := []int{1, 2, 3}

	i := 1
	for x := range Iter(xs) {
		if i != x {
			t.Errorf("Expected %d, got %d", i, x)
		}
		i++
	}
}
