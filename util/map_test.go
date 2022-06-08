package util

import (
	"reflect"
	"testing"
)

func TestGenMap(t *testing.T) {
	xs := []int{1, 2, 3}
	square := func(o int) int {
		return o * o
	}

	squares := make([]int, 0, 3)
	for x := range Map(Iter(xs), square) {
		squares = append(squares, x)
	}

	if expected := []int{1, 4, 9}; !reflect.DeepEqual(squares, expected) {
		t.Errorf("Expected %v, got %v", expected, squares)
	}
}
