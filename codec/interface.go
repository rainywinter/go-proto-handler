package codec

import "reflect"

type Message interface{}

// encode decode interface
type Codec interface {
	Encode(id int32, seq uint64, v Message, encodeErr error) ([]byte, error)
	Decode(b []byte, f func(id int32) (t reflect.Type, err error)) (res Message, id int32, seq uint64, err error)
}

// var codec Codec

// func SetCodec(c Codec) {
// 	codec = c
// }

// func Encode(id int32, seq uint64, v Message, encodeErr error) ([]byte, error) {
// 	return codec.Encode(id, seq, v, encodeErr)
// }

// func Decode(b []byte, f func(id int32, seq uint64) (t reflect.Type, err error)) (Message, error) {
// 	return codec.Decode(b, f)
// }
