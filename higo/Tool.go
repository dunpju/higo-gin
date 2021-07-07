package higo

import (
	"flag"
	"github.com/dengpju/higo-gin/higo/templates"
	"github.com/dengpju/higo-utils/utils"
	"log"
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

func (this *Tool) Cmd() {
	flag.StringVar(&this.Gen, "gen", "", "explain: Generate Controller or Model \n --option[controller | model] \n eg:-gen=controller")
	flag.StringVar(&this.Name, "name", "", "explain: Generate Name \neg:-name=Test")
	flag.StringVar(&this.Out, "out", "", "explain: Generate file output path \neg:-out=test\\app\\Controllers")
	//解析命令行参数
	flag.Parse()
	if "" != this.Gen {
		if controller == this.Gen {
			if this.Name == "" {
				log.Fatalln("name unable empty eg:-name=Test")
			}
			if this.Out == "" {
				log.Fatalln("out unable empty eg:-out=test\\app\\Controllers")
			}
			this.Package = utils.Basename(this.Out)
			//go run test\bin\main.go -gen=controller -name=Test -out=test\app\Controllers
			templates.NewController(this.Package, this.Name, this.Out).Generate()
		} else if model == this.Gen {

		} else {
			log.Fatalln("gen error option controller or model \neg:-gen=controller")
		}
		os.Exit(1)
	}
}
