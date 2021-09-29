package server

// base type for protocol messages, compatible with JSON (un)marshaling
type ProtMsg struct {
	Msg  string
	Data map[string]interface{}
}

// protocol message bound to a connection ID
type Msg struct {
	Id   int
	Data ProtMsg
}

// channels that a subscriber must listen on
type Subscriber struct {
	Join  chan int
	Leave chan int
	Rcv   chan Msg
}

func MakeSubscriber(buffer_size int) *Subscriber {
	return &Subscriber{
		Join:  make(chan int, buffer_size),
		Leave: make(chan int, buffer_size),
		Rcv:   make(chan Msg, buffer_size),
	}
}

type Server interface {
	Send(id int, msg string, data map[string]interface{}) error
	Subscribe(sub *Subscriber)
}
