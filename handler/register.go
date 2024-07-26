package handler

import (
	"context"
	"fmt"
	"reflect"

	"conn-proto-handler/codec"
)

type HandlerFunc func(context.Context, codec.Message) (codec.Message, error)

type MsgHandler struct {
	// store reflect.valueOf(*struct).type.elem(),
	// it will be used by reflect.New(req) later to create instance of req passed to handler
	ReqType reflect.Type
	// req msg handler
	Handler HandlerFunc
}

type MsgHandlerInfo struct {
	Handler      map[int32]MsgHandler
	ReqMsgTypeId map[reflect.Type]int32
	ResMsgIdType map[int32]reflect.Type
}

func (m *MsgHandlerInfo) GetMsgTypeById(id int32) (reflect.Type, error) {
	if codec.IsResMsg(id) {
		t, ok := m.ResMsgIdType[id]
		if !ok {
			return nil, fmt.Errorf("response msg %d not register", id)
		}
		return t, nil
	}
	h, ok := m.Handler[id]
	if !ok {
		return nil, fmt.Errorf("request msg %d not register", id)
	}
	return h.ReqType, nil
}

var GMsgHandlerInfo *MsgHandlerInfo

func init() {
	GMsgHandlerInfo = &MsgHandlerInfo{
		Handler:      make(map[int32]MsgHandler),
		ReqMsgTypeId: make(map[reflect.Type]int32),
		ResMsgIdType: make(map[int32]reflect.Type),
	}
}

// 约定主动发起消息id为奇数，回复的消息号为发送id+1，偶数
func RegisterMsg[T1 codec.Message, T2 codec.Message](msgId int32, req T1, res T2, cb func(context.Context, T1) (T2, error)) {
	var h HandlerFunc
	if cb != nil {
		h = func(ctx context.Context, v codec.Message) (codec.Message, error) {
			return cb(ctx, v.(T1))
		}
		GMsgHandlerInfo.Handler[msgId] = MsgHandler{
			ReqType: reflect.ValueOf(req).Type().Elem(),
			Handler: h,
		}
	} else {
		GMsgHandlerInfo.ReqMsgTypeId[reflect.ValueOf(req).Type()] = msgId
		GMsgHandlerInfo.ResMsgIdType[codec.ResMsgIdFromReq(msgId)] = reflect.ValueOf(res).Type().Elem()
	}
}
