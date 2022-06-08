package util

import (
	"errors"
	"golang.org/x/exp/constraints"
)

// Calculate the average of an iterable of floats
func Avg[T constraints.Float](iter Iterable[T]) (T, error) {
	var sum, count T
	for x := range iter {
		sum += x
		count += 1
	}
	if count == 0 {
		return 0, errors.New("Empty sequence")
	}
	return sum / count, nil
}
