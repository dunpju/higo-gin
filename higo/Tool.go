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
	if len(os.Args) >= 2 {
		flag.StringVar(&this.Gen, "gen", "", "explain: Generate Controller or Model \n --option[controller | model] \n eg:-gen=controller")
		flag.StringVar(&this.Name, "name", "", "explain: Generate Name \neg:-name=Test")
		flag.StringVar(&this.Out, "out", "", "explain: Generate file output path \neg:-out=test\\app\\Controllers")
		flag.Parse()
		if controller == this.Gen {
			if this.Name == "" {
				log.Fatalln("name unable empty \neg: -name=Test")
			}
			if this.Out == "" {
				log.Fatalln("out unable empty \neg: -out=test\\app\\Controllers")
			}
			this.Package = utils.Basename(this.Out)
			templates.NewController(this.Package, this.Name, this.Out).Generate()
		} else if model == this.Gen {
			if this.Name == "" {
				log.Fatalln("\ntable name unable empty \neg: -name=ts_user")
			}
			if this.Out == "" {
				log.Fatalln("\noutput directory unable empty \neg: -out=test\\app\\Models")
			}
			if this.Name == "all" {
				newModel := templates.NewModel(newGorm(), this.Name, this.Out, GetDbConfig().Database, GetDbConfig().Prefix)
				tables := newModel.GetTables()
				for _, table := range tables {
					templates.NewModel(newGorm(), table.Name, this.Out, GetDbConfig().Database, GetDbConfig().Prefix).Generate()
				}
			} else {
				templates.NewModel(newGorm(), this.Name, this.Out, GetDbConfig().Database, GetDbConfig().Prefix).Generate()
			}
		} else {
			log.Fatalln("\n gen Arguments Error! \nExplain: Generate Controller or Model \n --option[controller | model] \n eg:-gen=controller")
		}
		os.Exit(1)
	}
}
