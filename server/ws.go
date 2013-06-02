package server

import (
	"errors"
	"log"

	"code.google.com/p/go.net/websocket"
)

const (
	buffer_size = 64
)

type WsMsg struct {
	conn *Conn
	data ProtMsg
}

type WsServer struct {
	conns map[*Conn] int
	nextId int

	join chan *Conn
	leave chan *Conn
	rcv chan WsMsg
	send chan WsMsg

	subs []*Subscriber
}

func (s *WsServer) MakeWsHandler() websocket.Handler {
	handle := func(ws *websocket.Conn) {
		log.Printf("WsServer: New connection")
		conn := &Conn{ send: make (chan ProtMsg, 1024), ws: ws, server: s }
		s.join <- conn
		defer func() {
			log.Printf("WsServer: Lost connection")
			s.leave <- conn
		}()
		go conn.writer()
		conn.reader()
	}
	return websocket.Handler(handle)
}

func (s *WsServer) Subscribe(sub *Subscriber) {
	s.subs = append(s.subs, sub)
} 

func (s *WsServer) Run() {
	for {
		select {
		case conn := <-s.join:
			s.conns[conn] = s.nextId
			for _, sub := range s.subs {
				sub.Join <- s.nextId
			}
			s.nextId++
		case conn := <-s.leave:
			for _, sub := range s.subs {
				sub.Leave <- s.conns[conn]
			}
			delete(s.conns, conn)
			close(conn.send)
		case msg := <-s.rcv:
			for _, sub := range s.subs {
				sub.Rcv <- Msg{ Id: s.conns[msg.conn], Data: msg.data }
			}
		case msg := <-s.send:
			select {
			case msg.conn.send <- msg.data:
			default:
				log.Println("Error sending, closing connection")
				delete(s.conns, msg.conn)
				close(msg.conn.send)
				go msg.conn.ws.Close()
			}
		}
	}
}

func (s *WsServer) Send(id int, msg string, data map[string]interface{}) error {
	for conn, _id := range s.conns {
		if id == _id {
			conn.send <- ProtMsg{msg, data}
			return nil
		}
	}
	return errors.New("Connection not found")
}

type Conn struct {
	server *WsServer
	ws *websocket.Conn
	send chan ProtMsg
}

func (c *Conn) reader() {
	for {
		var data ProtMsg
		err := websocket.JSON.Receive(c.ws, &data)
		if err != nil {
			log.Println(err)
			break
		}
		var msg = WsMsg{ conn: c, data: data }
		c.server.rcv <- msg
	}
	c.ws.Close()
}

func (c *Conn) writer() {
	for message := range c.send {
		err := websocket.JSON.Send(c.ws, message)
		if err != nil {
			log.Println(err)
			break
		}
	}
	c.ws.Close()
}

func NewServer() *WsServer {
	return &WsServer{
		conns: make(map[*Conn]int, buffer_size),
		join: make(chan *Conn, buffer_size),
		leave: make(chan *Conn, buffer_size),
		rcv: make(chan WsMsg, buffer_size),
		send: make(chan WsMsg, buffer_size),
	}
}
