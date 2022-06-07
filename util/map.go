package util

type MapFunc[T any, U any] func(T) U

// Lazily map a function to an Iterable. Returns a new Iterable
func Map[T any, U any](xi Iterable[T], mapFunc MapFunc[T, U]) Iterable[U] {
	ch := make(chan U)

	go (func(out chan<- U) {
		for x := range xi {
			out <- mapFunc(x)
		}
		close(out)
	})(ch)

	return ch
}
