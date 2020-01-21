package socketgo

type Message struct {
	Event   string
	Payload interface{}
}

func NewMessage(event string, payload interface{}) *Message {
	return &Message{
		Event:   event,
		Payload: payload,
	}
}
