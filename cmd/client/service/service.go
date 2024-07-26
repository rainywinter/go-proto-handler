package service

import (
	pb "conn-proto-handler/cmd/protobuf"
	"conn-proto-handler/handler"
	"conn-proto-handler/server"
	"context"
	"fmt"
)

type Service struct {
	id   int32
	conn *server.Conn
}

func NewService(id int32, conn *server.Conn) *Service {
	s := &Service{
		id:   id,
		conn: conn,
	}
	s.register()

	return s
}

func (s *Service) register() {
	handler.RegisterMsg(int32(pb.MsgId_Login), &pb.LoginRq{}, &pb.LoginRs{}, nil)
	handler.RegisterMsg(int32(pb.MsgId_Hello), &pb.Hello_Rq{}, &pb.Hello_Rs{}, s.Hello)
}

func (s *Service) Hello(ctx context.Context, req *pb.Hello_Rq) (*pb.Hello_Rs, error) {
	fmt.Println("Hello:req", req.Name)
	return &pb.Hello_Rs{Msg: fmt.Sprintf("hello, i am client:%d", s.id)}, nil
}

func (s *Service) Login(ctx context.Context) {
	res, err := s.conn.SendMessage(ctx, &pb.LoginRq{Id: s.id})
	fmt.Printf("client:%d send login, res:%v, err:%v)\n", s.id, res, err)
}
