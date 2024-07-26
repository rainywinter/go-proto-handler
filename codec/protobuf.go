package codec

import (
	pb "conn-proto-handler/cmd/protobuf"
	"errors"
	"reflect"

	"google.golang.org/protobuf/proto"
)

type protoCodec struct {
}

func NewProtoCodec() Codec {
	return &protoCodec{}
}

func (p *protoCodec) Encode(id int32, seq uint64, v Message, encodeErr error) ([]byte, error) {
	m := &pb.Message{
		Id:  pb.MsgId(id),
		Seq: seq,
	}
	if encodeErr == nil {
		body, err := proto.Marshal(v.(proto.Message))
		if err != nil {
			return nil, err
		}
		m.Body = body
	} else {
		m.ErrMsg = encodeErr.Error()
	}

	buf, err := proto.Marshal(m)
	return buf, err
}
func (p *protoCodec) Decode(b []byte, f func(id int32) (t reflect.Type, err error)) (res Message, id int32, seq uint64, err error) {
	m := &pb.Message{}
	err = proto.Unmarshal(b, m)
	if err != nil {
		return
	}
	id = int32(m.Id)
	seq = m.Seq

	if m.ErrMsg != "" {
		res = errors.New(m.ErrMsg)
		return
	}

	t, err := f(int32(m.Id))
	if err != nil {
		return
	}
	v := reflect.New(t).Interface()
	err = proto.Unmarshal(m.Body, v.(proto.Message))
	if err != nil {
		return
	}

	res = v
	return
}
