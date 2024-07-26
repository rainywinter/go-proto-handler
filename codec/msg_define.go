package codec

func ResMsgIdFromReq(reqId int32) int32 { return reqId + 1 }
func IsResMsg(id int32) bool            { return id&1 == 0 }
