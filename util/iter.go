package util

// General interface for iterable collections. Implements iterators as a channel. Due to lack of generics, we use interface{}, and need to cast back up to the actual type.
type Iterable interface {
	Iter() <-chan interface{}
}

// An iterable integer slice
type IntSlice []int

func (xs IntSlice) Iter() <-chan interface{} {
	ch := make(chan interface{})

	go (func() {
		for _, x := range xs {
			ch <- x
		}
		close(ch)
	})()

	return ch
}

// An iterable float32 slice
type FloatSlice []float32

func (xs FloatSlice) Iter() <-chan interface{} {
	ch := make(chan interface{})

	go (func() {
		for _, x := range xs {
			ch <- x
		}
		close(ch)
	})()

	return ch
}

type mapFunc func(interface{}) interface{}

type mappedIterable struct {
	f    mapFunc
	iter Iterable
}

func (m *mappedIterable) Iter() <-chan interface{} {
	ch := make(chan interface{})

	go (func(out chan<- interface{}) {
		for x := range m.iter.Iter() {
			out <- m.f(x)
		}
		close(out)
	})(ch)

	return ch
}

// Lazily map a function to an Iterable. Returns a new Iterable.
func Map(f mapFunc, iter Iterable) Iterable {
	return &mappedIterable{f, iter}
}
