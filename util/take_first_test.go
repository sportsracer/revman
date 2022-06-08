package util

import (
	"fmt"
	"reflect"
	"testing"
)

func (xi Iterable[T]) equalsArray(expected []T) error {
	count := len(expected)
	xs := make([]T, 0, count)
	for x := range xi {
		xs = append(xs, x)
	}

	if !reflect.DeepEqual(xs, expected) {
		return fmt.Errorf("Expected %v, got %v", expected, xs)
	}
	return nil
}

func TestTakeFirst(t *testing.T) {
	t.Run("takes the first n elements", func(t *testing.T) {
		xs := []int{1, 2, 3}

		sliced := TakeFirst(Iter(xs), 2)

		if err := sliced.equalsArray([]int{1, 2}); err != nil {
			t.Error(err)
		}
	})

	t.Run("works on inifite iterables", func(t *testing.T) {
		yes := func() Iterable[string] {
			ch := make(chan string)
			go (func() {
				for {
					ch <- "y"
				}
			})()
			return ch
		}

		sliced := TakeFirst(yes(), 3)

		if err := sliced.equalsArray([]string{"y", "y", "y"}); err != nil {
			t.Error(err)
		}
	})
}
