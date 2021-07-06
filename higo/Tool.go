package higo

import (
	"flag"
	"github.com/dengpju/higo-gin/higo/templates"
	"github.com/dengpju/higo-utils/utils"
	"os"
)

const (
	controller = "controller"
	model      = "model"
)

type Tool struct {
	Gen     string
	Name    string
	Out     string
	Package string
}

func NewTool() *Tool {
	return &Tool{}
}

func (this *Tool) Execute() {
	flag.StringVar(&this.Gen, "gen", "", "explain: Generate Controller or Model \n --option[controller | model]")
	flag.StringVar(&this.Name, "name", "", "explain: Generate Name")
	flag.StringVar(&this.Out, "out", "", "explain: Generate file output path")
	//解析命令行参数
	flag.Parse()
	if "" != this.Gen {
		if controller == this.Gen {
			this.Package = utils.Basename(this.Out)
			//go run test\bin\main.go -gen=controller -name=Test -out=test\app\Controllers
			templates.NewController(this.Package, this.Name, this.Out).Generate()
		} else if model == this.Gen {

		} else {
			panic("gen error; option controller or model")
		}
		os.Exit(1)
	}
}
