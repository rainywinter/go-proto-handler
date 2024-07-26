package service

import (
	pb "conn-proto-handler/cmd/protobuf"
	"conn-proto-handler/codec"
	"conn-proto-handler/delivery"
	"conn-proto-handler/handler"
	"conn-proto-handler/server"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Service struct {
	upgrader    websocket.Upgrader
	connManager *server.ConnManager
}

func NewService(upgrader websocket.Upgrader) *Service {
	s := &Service{
		upgrader:    upgrader,
		connManager: server.NewConnManager(),
	}
	s.register()

	return s
}

func (s *Service) register() {
	handler.RegisterMsg(int32(pb.MsgId_Login), &pb.LoginRq{}, &pb.LoginRs{}, s.Login)
	// handler.RegisterMsg(int32(pb.MsgId_Hello), &pb.Hello_Rq{}, &pb.Hello_Rs{}, s.Hello)
	handler.RegisterMsg(int32(pb.MsgId_Hello), &pb.Hello_Rq{}, &pb.Hello_Rs{}, nil)
}

// serve http upgrade
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsConn, err := s.upgrader.Upgrade(w, r, http.Header{})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println("new conn:", r.RemoteAddr)
	conn := server.NewConn(codec.NewProtoCodec(), delivery.NewWebsocketConn(wsConn), s.connManager.UnexpectClose)
	conn.Run()
	// conn exit
}

func (s *Service) Shutdown() {
	s.connManager.Clear()
}

func (s *Service) Login(ctx context.Context, req *pb.LoginRq) (*pb.LoginRs, error) {
	fmt.Println("Login:req", req.Id)

	conn, ok := ctx.Value(server.ConnCtxKey{}).(*server.Conn)
	if !ok {
		return nil, errors.New("login err, not conn type")
	}
	id := server.ConnID(req.Id)
	conn.SetConnId(server.ConnID(req.Id))
	s.connManager.Add(id, conn)

	return &pb.LoginRs{Ok: true}, nil
}

func (s *Service) Hello(ctx context.Context, id server.ConnID, hello string) {
	fmt.Println("Hello send:", "id", id, "content", hello)

	conn := s.connManager.Get(id)
	if conn == nil {
		fmt.Println("Hello:conn not logined", "id", id)
		return
	}
	res, err := conn.SendMessage(ctx, &pb.Hello_Rq{Msg: hello})
	fmt.Println("Hello:res", "id", id, "content", hello, "res", res, "err", err)
}
