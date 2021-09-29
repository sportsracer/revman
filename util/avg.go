package util

import (
	"errors"
)

// Calculate the average of an iterable of float32, e.g. a FloatSlice
func Avg(iter Iterable) (float32, error) {
	var sum float32
	var count int
	for x := range iter.Iter() {
		sum += x.(float32)
		count += 1
	}
	if count == 0 {
		return 0, errors.New("Empty sequence")
	}
	return sum / float32(count), nil
}
