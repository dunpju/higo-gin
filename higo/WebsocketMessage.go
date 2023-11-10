package higo

import (
	"encoding/json"
	"github.com/dunpju/higo-utils/utils/maputil"
)

type WsReadMessage struct {
	MessageType int
	MessageData []byte
}

func NewReadMessage(messageType int, messageData []byte) *WsReadMessage {
	return &WsReadMessage{MessageType: messageType, MessageData: messageData}
}

type WsWriteMessage struct {
	MessageType string
	MessageData []byte
}

func WsRespString(messageData string) WsWriteMessage {
	return WsWriteMessage{MessageType: WsRespstring, MessageData: []byte(messageData)}
}

func WsRespMap(messageData *maputil.ArrayMap) WsWriteMessage {
	return WsWriteMessage{MessageType: WsRespmap, MessageData: []byte(messageData.String())}
}

func WsRespStruct(messageData interface{}) WsWriteMessage {
	mjson, err := json.Marshal(messageData)
	if err != nil {
		panic(err)
	}
	return WsWriteMessage{MessageType: WsRespstruct, MessageData: mjson}
}

func WsRespError(messageData string) WsWriteMessage {
	return WsWriteMessage{MessageType: WsResperror, MessageData: []byte(messageData)}
}

func WsRespClose() WsWriteMessage {
	return WsWriteMessage{MessageType: WsRespclose, MessageData: []byte("1")}
}
