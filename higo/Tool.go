package higo

import (
	"flag"
	"fmt"
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
		var tplEngine templates.ItplEngine
		if controller == this.Gen {
			if this.Name == "" {
				log.Fatalln("name unable empty eg: -name=Test")
			}
			if this.Out == "" {
				log.Fatalln("out unable empty eg: -out=test\\app\\Controllers")
			}
			this.Package = utils.Basename(this.Out)
			tplEngine = templates.NewController(this.Package, this.Name, this.Out)
		} else if model == this.Gen {
			if this.Name == "" {
				log.Fatalln("name unable empty eg: -name=ts_user")
			}
			if this.Out == "" {
				log.Fatalln("out unable empty eg: -out=test\\app\\Models")
			}
			tplEngine = templates.NewModel(newGorm(), this.Name, this.Out, GetDbConfig().Database, GetDbConfig().Prefix)
			fmt.Println(tplEngine)
		} else {
			log.Fatalln("gen error option controller or model \neg:-gen=controller")
		}
		tplEngine.Generate()
		os.Exit(1)
	}
}
