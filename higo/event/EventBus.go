package event

import "sync"

type EventBus struct {
	subscribes map[string]EventDataChannel
	lock       sync.RWMutex //读写锁
}

func NewEventBus() *EventBus {
	return &EventBus{subscribes: make(map[string]EventDataChannel)}
}

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
