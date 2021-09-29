package util

import "testing"

func TestAvg(t *testing.T) {
	t.Run("returns an error on empty iterables", func(t *testing.T) {
		xs := FloatSlice{}

		_, err := Avg(xs)

		if err == nil {
			t.Error("Expected error here")
		}
	})

	t.Run("computes average correctly", func(t *testing.T) {
		xs := FloatSlice{1.0, 2.0, 3.0}

		avg, err := Avg(xs)

		if err != nil {
			t.Error("Unexpected error")
		}

		if expected := float32(2.0); avg != expected {
			t.Errorf("Expected %f, got %f", expected, avg)
		}
	})
}
