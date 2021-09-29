package util

import "testing"

func TestMax(t *testing.T) {
	t.Run("returns nil for empty iterables", func(t *testing.T) {
		xs := IntSlice{}

		maxItem, max := MaxOccur(xs)

		if expectedMaxItem := interface{}(nil); maxItem != expectedMaxItem {
			t.Errorf("Expected %d, got %d", expectedMaxItem, maxItem)
		}

		if expectedMax := 0; max != expectedMax {
			t.Errorf("Expected %d, got %d", expectedMax, max)
		}
	})

	t.Run("identifies the most frequent element in an iterable", func(t *testing.T) {
		xs := IntSlice{1, 1, 2, 3, 5}

		maxItem, max := MaxOccur(xs)

		if expectedMaxItem := 1; maxItem != expectedMaxItem {
			t.Errorf("Expected %d, got %d", expectedMaxItem, maxItem)
		}

		if expectedMax := 2; max != expectedMax {
			t.Errorf("Expected %d, got %d", expectedMax, max)
		}
	})
}
