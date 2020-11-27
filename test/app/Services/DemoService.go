package Services

type DemoService struct {
	Demo string
}

func NewDemoService() *DemoService {
	return &DemoService{Demo: "demo"}
}

