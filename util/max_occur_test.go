package util

import "testing"

type Empty struct{}

func TestMax(t *testing.T) {
	t.Run("returns nil for empty iterables", func(t *testing.T) {
		xs := []Empty{}

		val, count := MaxOccur(Iter(xs))

		if expectedVal := (*Empty)(nil); val != expectedVal {
			t.Errorf("Expected %d, got %d", expectedVal, *val)
		}

		if expectedCount := uint(0); count != expectedCount {
			t.Errorf("Expected %d, got %d", expectedCount, count)
		}
	})

	t.Run("identifies the most frequent element in an iterable", func(t *testing.T) {
		xs := []int{1, 1, 2, 3, 5}

		val, count := MaxOccur(Iter(xs))

		if expectedVal := 1; *val != expectedVal {
			t.Errorf("Expected %d, got %d", expectedVal, *val)
		}

		if expectedCount := uint(2); count != expectedCount {
			t.Errorf("Expected %d, got %d", expectedCount, count)
		}
	})
}
