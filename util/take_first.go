package util

// Take the first `count` elements from the beginning of an iterable
func TakeFirst[T any](xi Iterable[T], count int) Iterable[T] {
	ch := make(chan T)

	go (func(out chan<- T) {
		i := 0
		for x := range xi {
			if i++; i > count {
				break
			}
			out <- x
		}
		close(out)
	})(ch)

	return ch
}
