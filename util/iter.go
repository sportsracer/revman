package util

// Iterables are implemented as Go channels
type Iterable[T any] <-chan T

// Convert an array to an iterable by sending its elements to a channel
func Iter[T any](xs []T) Iterable[T] {
	ch := make(chan T)

	go (func() {
		for _, x := range xs {
			ch <- x
		}
		close(ch)
	})()

	return ch
}
