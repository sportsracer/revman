package util

import (
	"reflect"
	"testing"
)

func TestIterable(t *testing.T) {
	xs := IntSlice{1, 2, 3}

	i := 1
	for x := range xs.Iter() {
		if i != x {
			t.Errorf("Expected %d, got %d", i, x)
		}
		i++
	}
}

func TestMap(t *testing.T) {
	xs := IntSlice{1, 2, 3}
	squares := make([]int, 0, 3)
	square := func(o interface{}) interface{} {
		return o.(int) * o.(int)
	}

	for x := range Map(square, xs).Iter() {
		squares = append(squares, x.(int))
	}

	if expected := []int{1, 4, 9}; !reflect.DeepEqual(squares, expected) {
		t.Errorf("Expected %v, got %v", expected, squares)
	}
}
