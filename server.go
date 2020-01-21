package socketgo

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
)

type clientMap map[net.Addr]*Conn
type handlerFunc func(conn *Conn, data interface{}) error
type handlerMap map[string]handlerFunc

type Server struct {
	onConnect func(conn *Conn) error
	clients   clientMap
	handlers  handlerMap
	upgrader  *websocket.Upgrader
}

func (s *Server) Broadcast(event string, payload interface{}) error {
	encoded, err := json.Marshal(NewMessage(event, payload))
	if err != nil {
		return err
	}

	for addr, conn := range s.clients {
		if err := conn.SendBytes(encoded); err != nil {
			log.Printf("Failed to write message to %s: %s\n", addr.String(), err)
		}
	}

	return nil
}

func (s *Server) Handle(event string, handler handlerFunc) {
	s.handlers[event] = handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection for %s: %s\n", r.RemoteAddr, err)
		return
	}

	s.register(NewConn(conn))
	s.handleConn(NewConn(conn))
}

func (s *Server) OnConnect(fn func(*Conn) error) {
	s.onConnect = fn
}

func (s *Server) handleConn(conn *Conn) {
	defer s.closeConn(conn)

	if s.onConnect != nil {
		if err := s.onConnect(conn); err != nil {
			log.Printf("Error executing onConnect for %s: %s\n", conn.RemoteAddr(), err)
		}
	}

	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message for %s: %s\n", conn.RemoteAddr(), err)
			break
		}

		if messageType != websocket.TextMessage {
			continue
		}

		message := &Message{}
		if err := json.Unmarshal(payload, message); err != nil {
			log.Printf("Error decoding payload %s for %s: %s\n", payload, conn.RemoteAddr(), err)
			continue
		}

		if err := s.handleMessage(conn, message); err != nil {
			log.Printf("Error handling message %+v for %s: %s\n", message, conn.RemoteAddr(), err)
			continue
		}
	}
}

func (s *Server) handleMessage(conn *Conn, message *Message) error {
	handler, ok := s.handlers[message.Event]
	if !ok {
		return nil
	}

	return handler(conn, message.Payload)
}

func (s *Server) closeConn(conn *Conn) {
	if err := conn.Close(); err != nil {
		log.Printf("Error closing conn for %s: %s", conn.RemoteAddr(), err)
	}

	s.unregister(conn)
}

func (s *Server) register(conn *Conn) {
	s.clients[conn.RemoteAddr()] = conn
}

func (s *Server) unregister(conn *Conn) {
	delete(s.clients, conn.RemoteAddr())
}

func NewServer() *Server {
	return &Server{
		clients:  clientMap{},
		handlers: handlerMap{},
		upgrader: &websocket.Upgrader{},
	}
}
