package tests

import (
	pb "conn-proto-handler/cmd/protobuf"
	"conn-proto-handler/handler"
	"context"
	"fmt"
	"testing"

	"reflect"
)

type Service struct{}

func (s *Service) Hello(ctx context.Context, req *pb.Hello_Rq) (*pb.Hello_Rs, error) {
	fmt.Println("Hello:req", req.Name)
	return &pb.Hello_Rs{Msg: "hello middleware"}, nil
}

func register() {
	s := &Service{}
	handler.RegisterMsg(int32(pb.MsgId_Hello), &pb.Hello_Rq{}, &pb.Hello_Rs{}, s.Hello)
	handler.RegisterMsg(int32(pb.MsgId_Hello), &pb.Hello_Rq{}, &pb.Hello_Rs{}, nil)
}

func TestWithoutConn(t *testing.T) {
	register()

	buf, err := pb.Encode(int32(pb.MsgId_Hello), 1, &pb.Hello_Rq{Name: "testserver1"}, nil)
	if err != nil {
		panic(err)
	}

	var msgId int32
	msg, err := pb.Decode(buf, func(id int32, seq uint64) (t reflect.Type, err error) {
		msgId = id
		if id&1 == 0 {
			return handler.GMsgHandlerInfo.ResMsgIdType[id], nil
		}
		return handler.GMsgHandlerInfo.Handler[id].ReqType, nil
	})
	if err != nil {
		panic(err)
	}

	h := handler.GMsgHandlerInfo.Handler[msgId]
	res, err := h.Handler(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

}
