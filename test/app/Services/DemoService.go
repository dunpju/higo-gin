package Services

import "github.com/dengpju/higo-gin/higo/event"

const GetDemoList  = "GetList"
var Bus *event.EventBus
var DemoListCh event.EventDataChannel

func init() {
	Bus = event.NewEventBus()
}

type DemoService struct {
	Demo string
}

func GetDemoListCh() event.EventDataChannel {
	return Bus.Sub(GetDemoList, NewDemoService().GetList)
}

func NewDemoService() *DemoService {
	return &DemoService{Demo: "demo"}
}

func (this *DemoService) GetList() []interface{} {
	return []interface{}{
		struct {
			Id   int
			Name string
		}{Id: 101, Name: "aaa"},
		struct {
			Id   int
			Name string
		}{Id: 102, Name: "bbb"},
	}
}
