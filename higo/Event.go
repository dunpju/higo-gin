package higo

import (
	"sync"
)

func init() {
	events = newEventRepository()
}

var (
	events *EventRepository
)

const (
	// BeforeStart Before server start
	BeforeStart EventType = iota + 1
	// AfterLoadRoute After Load Route
	AfterLoadRoute
)

func eventRegister(iEvent IEvent) {
	events.add(iEvent.Event(), iEvent)
}

func eventPoint(hg *Higo, event EventType) {
	e := events.get(event)
	if e != nil {
		e.Handle(hg)
	}
}

type IEvent interface {
	Event() EventType
	Handle(hg *Higo)
}

type EventType int

type EventRepository struct {
	sort    []EventType
	syncMap sync.Map
}

func newEventRepository() *EventRepository {
	return &EventRepository{sort: []EventType{}}
}

func (this *EventRepository) Len() int {
	return len(this.sort)
}

func (this *EventRepository) get(event EventType) IEvent {
	e, ok := this.syncMap.Load(event)
	if ok {
		return e.(IEvent)
	}
	return nil
}

func (this *EventRepository) add(event EventType, iEvent IEvent) {
	_, ok := this.syncMap.Load(event)
	if !ok {
		this.sort = append(this.sort, event)
	}
	this.syncMap.Store(event, iEvent)
}

type EventHandle func(hg *Higo)

type Event struct {
	eventType EventType
	handle    EventHandle
}

func addEvent(eventType EventType, handle EventHandle) {
	eventRegister(&Event{eventType: eventType, handle: handle})
}

func (this *Event) Event() EventType {
	return this.eventType
}

func (this *Event) Handle(hg *Higo) {
	this.handle(hg)
}
