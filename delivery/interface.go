package delivery

type Delivery interface {
	ReadMessage() ([]byte, error)
	WriteMessage([]byte) error
	Close() error
}
