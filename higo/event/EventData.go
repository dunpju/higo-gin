package event

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

type EventData struct {
	Data interface{}
}

func NewEventData(data interface{}) *EventData {
	return &EventData{Data: data}
}

type EventDataChannel chan *EventData

func (this EventDataChannel) Data(timeout time.Duration) interface{} {
	ctx, cancle := context.WithTimeout(context.Background(), timeout)
	defer cancle()
	select {
	case <-ctx.Done(): //超时
		return gin.H{"message": "timeout"}
	case data := <-this:
		return data.Data
	}
}
