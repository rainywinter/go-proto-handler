package delivery_test

import (
	"conn-proto-handler/delivery"
	"testing"
	"time"
)

func TestChannelConn(t *testing.T) {
	t.Log("==>testChanConn")
	// serverSide
	recvCh, sendCh := make(chan []byte, 10), make(chan []byte, 10)
	serverIO := delivery.NewChanConn(recvCh, sendCh)
	clientIO := delivery.NewChanConn(sendCh, recvCh)

	go func() {
		serverIO.WriteMessage([]byte("hello, i am server"))

		msg, err := clientIO.ReadMessage()
		t.Log(string(msg), err)
		serverIO.Close()

	}()

	clientIO.WriteMessage([]byte("hello, i am client"))

	msg, err := serverIO.ReadMessage()
	t.Log(string(msg), err)

	clientIO.Close()
	time.Sleep(time.Second)
}
