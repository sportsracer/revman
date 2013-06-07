package util

type Iterable interface {
	Iter() <-chan interface{}
}

// make int slices iterable

type IntSlice []int

func (xs IntSlice) iterate(out chan<- interface{}) {
	for x := range xs {
		out <- x
	}
	close(out)
}

func (xs IntSlice) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go xs.iterate(ch)
	return ch
}

// map
type mapFunc func(interface{}) interface{}

type mappedIterable struct {
	f mapFunc
	iter Iterable
}

func (m *mappedIterable) iterate(out chan<- interface{}) {
	for x := range m.iter.Iter() {
		out <- m.f(x)
	}
	close(out)
}

func (m *mappedIterable) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go m.iterate(ch)
	return ch
}

func Map(f mapFunc, iter Iterable) Iterable {
	return &mappedIterable{f, iter}
}
