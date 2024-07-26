package server

import (
	"conn-proto-handler/codec"
	"conn-proto-handler/delivery"
	"conn-proto-handler/handler"
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type ConnID int

type waitResCh chan interface{}
type ConnCtxKey struct{}
type Conn struct {
	id         ConnID
	logined    bool
	seq        uint64
	recvQ      chan []byte
	sendQ      chan []byte
	coder      codec.Codec
	express    delivery.Delivery
	exitFunc   func(id ConnID) // 异常退出回调
	waitResSet sync.Map        // key is msg id, value is waitResCh
	closeOnce  sync.Once
	closeCh    chan struct{}
}

func NewConn(coder codec.Codec, express delivery.Delivery, exitFunc func(id ConnID)) *Conn {
	return &Conn{
		seq:      uint64(time.Now().Unix()),
		recvQ:    make(chan []byte, 10),
		sendQ:    make(chan []byte, 10),
		coder:    coder,
		express:  express,
		exitFunc: exitFunc,
		closeCh:  make(chan struct{}),
	}
}

func (c *Conn) SetConnId(id ConnID) {
	c.id = id
	c.logined = true
}

func (c *Conn) Run() {
	defer func() {
		_ = c.express.Close()
		close(c.closeCh)

		if c.exitFunc != nil {
			c.exitFunc(c.id)
		}
	}()
	go c.recv()
	go c.write()
	c.read()
}

func (c *Conn) Close() error {
	var err error
	c.closeOnce.Do(func() {
		c.exitFunc = nil
		err = c.express.Close()
	})

	return err
}

func (c *Conn) read() {
	defer func() {
		close(c.recvQ)
		fmt.Println("exit read")
	}()

	for {
		buf, err := c.express.ReadMessage()
		if err != nil {
			// todo
			fmt.Println("read message err", err)
			return
		}

		c.recvQ <- buf
	}
}

func (c *Conn) write() {
	for {
		select {
		case <-c.closeCh:
			fmt.Println("exit write")
			return
		case b := <-c.sendQ:
			err := c.express.WriteMessage(b)
			if err != nil {
				// todo
				fmt.Println("write message error", err)
			}
		}
	}
}

func (c *Conn) send(b []byte) error {
	select {
	case <-c.closeCh:
		return errors.New("channel closed, don't send")
	case c.sendQ <- b:
		return nil
	}
}

func (c *Conn) SendMessage(ctx context.Context, req codec.Message) (res codec.Message, err error) {
	seq := atomic.AddUint64(&c.seq, 1)

	msgId := handler.GMsgHandlerInfo.ReqMsgTypeId[reflect.ValueOf(req).Type()]
	fmt.Println("conn send:", msgId, "seq", seq)

	buf, err := c.coder.Encode(msgId, seq, req, nil)
	if err != nil {
		return
	}

	defer c.waitResSet.Delete(seq)
	ch := make(waitResCh, 1)
	c.waitResSet.Store(seq, ch)
	if err = c.send(buf); err != nil {
		return
	}

	if _, ok := handler.GMsgHandlerInfo.ResMsgIdType[codec.ResMsgIdFromReq(msgId)]; ok {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout")
		case msg := <-ch:
			if e, ok := msg.(error); ok {
				return nil, e
			}
			res = msg.(codec.Message)
			return
		}
	}
	return
}

func (c *Conn) recv() {
	defer func() {
		fmt.Println("exit recv")
	}()

	for dataBuf := range c.recvQ {
		ctx := context.Background()
		msg, msgId, seq, err := c.coder.Decode(dataBuf, handler.GMsgHandlerInfo.GetMsgTypeById)
		if err != nil {
			fmt.Println("read  msg err", err)
			continue
		}
		fmt.Println("recv msg:", msg, msgId, seq, err, codec.IsResMsg(msgId))
		if codec.IsResMsg(msgId) {
			// receive response
			if v, ok := c.waitResSet.LoadAndDelete(seq); ok {
				ch := v.(waitResCh)
				ch <- msg
			}
		} else {
			ctx = context.WithValue(ctx, ConnCtxKey{}, c)

			h := handler.GMsgHandlerInfo.Handler[msgId]
			res, err := h.Handler(ctx, msg)
			if err != nil {
				fmt.Println("handler err:", err)
			}

			// send response
			// err = errors.New("test respone err")
			buf, _ := c.coder.Encode(codec.ResMsgIdFromReq(msgId), seq, res, err)

			c.send(buf)
		}
	}
}
