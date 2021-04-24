package higo

type WebsocketMessage struct {
	MessageType int
	MessageData []byte
}

func NewWebsocketMessage(messageType int, messageData []byte) *WebsocketMessage {
	return &WebsocketMessage{MessageType: messageType, MessageData: messageData}
}


