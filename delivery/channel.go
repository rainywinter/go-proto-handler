package delivery

import (
	"errors"
	"sync/atomic"
)

type chanConn struct {
	recvCh <-chan []byte
	sendCh chan<- []byte
	closed uint32
}

func NewChanConn(recvCh chan []byte, sendCh chan []byte) Delivery {
	return &chanConn{
		recvCh: recvCh,
		sendCh: sendCh,
	}
}

func (c *chanConn) ReadMessage() ([]byte, error) {
	b, ok := <-c.recvCh
	if !ok {
		return nil, errors.New("read: channel closed")
	}
	return b, nil
}

func (c *chanConn) WriteMessage(b []byte) error {
	if atomic.CompareAndSwapUint32(&c.closed, 0, 0) {
		c.sendCh <- b
		return nil
	}
	return errors.New("send closed")
}

func (c *chanConn) Close() error {
	if atomic.CompareAndSwapUint32(&c.closed, 0, 1) {
		close(c.sendCh)
	}

	return nil
}
