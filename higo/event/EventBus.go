package event

import "sync"

type EventBus struct {
	subscribes map[string]EventDataChannel
	handlers   map[string]*EventHandler
	lock       sync.RWMutex //读写锁
}

func NewEventBus() *EventBus {
	return &EventBus{subscribes: make(map[string]EventDataChannel), handlers: make(map[string]*EventHandler)}
}

func (this *EventBus) Sub(topic string, fn interface{}) EventDataChannel {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.subscribes[topic] = make(EventDataChannel)
	this.handlers[topic] = NewEventHandler(fn)
	return this.subscribes[topic]
}

func (this *EventBus) Pub(topic string, arguments ...interface{}) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if ec, found := this.subscribes[topic]; found {
		handler := this.handlers[topic]
		go func() {
			ec <- &EventData{Data: handler.Call(arguments...)}
		}()
	}
}

/**
//订阅
func (this *EventBus) Sub(topic string) EventDataChannel {
	this.lock.Lock() //写锁
	defer this.lock.Unlock()
	if ec, found := this.subscribes[topic]; found {
		return ec
	} else {
		this.subscribes[topic] = make(EventDataChannel)
		return this.subscribes[topic]
	}
}

func (this *EventBus) Pub(topic string, data interface{}) {
	this.lock.RLock() //读锁
	defer this.lock.RUnlock()
	if ec, found := this.subscribes[topic]; found {
		go func() {
			ec <- NewEventData(data)
		}()
	}
}

*/
