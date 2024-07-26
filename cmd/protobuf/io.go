package message_pb

import (
	"errors"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// 该部分逻辑已经迁移到codec package下

func Encode(id int32, seq uint64, v interface{}, encodeErr error) ([]byte, error) {
	m := &Message{
		Id:  MsgId(id),
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

func Decode(b []byte, f func(id int32, seq uint64) (t reflect.Type, err error)) (interface{}, error) {
	m := &Message{}
	err := proto.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}

	if m.ErrMsg != "" {
		return nil, errors.New(m.ErrMsg)
	}

	t, err := f(int32(m.Id), m.Seq)
	if err != nil {
		return nil, err
	}
	v := reflect.New(t).Interface()
	err = proto.Unmarshal(m.Body, v.(proto.Message))
	if err != nil {
		return nil, err
	}

	return v, nil
}
