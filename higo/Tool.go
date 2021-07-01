package higo

import (
	"flag"
	"fmt"
	"github.com/dengpju/higo-gin/higo/templates"
	"github.com/dengpju/higo-utils/utils"
	"os"
)

const (
	controller = "controller"
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
	flag.StringVar(&this.Gen, "gen", "", "generate")
	flag.StringVar(&this.Name, "name", "", "name")
	flag.StringVar(&this.Out, "out", "", "out path")
	//解析命令行参数
	flag.Parse()
	if "" != this.Gen {
		if controller == this.Gen {
			fmt.Println(this.Gen)
			fmt.Println(this.Name)
			fmt.Println(this.Out)
			this.Package = utils.Basename(this.Out)
			fmt.Println(utils.Basename(this.Out))
			//go run test\bin\main.go -gen=controller -name=Test -out=test\app\Controllers
			templates.NewController(this.Package, this.Name, this.Out+"\\"+this.Name+"Controller.go").Generate()
		}
		os.Exit(1)
	}
}
