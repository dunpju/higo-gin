package {{.Package}}

import (
	"github.com/dengpju/higo-gin/higo"
)

type {{.Name}}Controller struct {
}

func New{{.Name}}Controller() *{{.Name}}Controller {
	return &{{.Name}}Controller{}
}

func (this *{{.Name}}Controller) New() higo.IClass {
	return New{{.Name}}Controller()
}

func (this *{{.Name}}Controller) Route(hg *higo.Higo) {
    //TODO::example
	/**
	//route example
	hg.Get("/relative", this.Example, hg.Flag("unique flag"), hg.Desc("description"))

    //route group example
    hg.AddGroup("/group_prefix", func() {
    	hg.Get("/relative", this.Example, hg.Flag("unique flag"), hg.Desc("description"))
    })
    */
}