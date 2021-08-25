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
	enum       = "enum"
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
		flag.StringVar(&this.Gen, "gen", "", `explain: Generate Controller or Model or Enum
--option[controller | model | enum]
eg:-gen=controller`)
		flag.StringVar(&this.Name, "name", "", `explain: Generate Name 
eg:-name=Test`)
		flag.StringVar(&this.Out, "out", "", `explain: Generate file output path 
eg:-out=test\app\Controllers`)
		flag.Parse()
		if controller == this.Gen {
			if this.Name == "" {
				log.Fatalln(`controller name unable empty 
eg: -name=Test`)
			}
			if this.Out == "" {
				log.Fatalln(`output directory unable empty 
eg: -out=test\app\Controllers`)
			}
			this.Package = utils.Basename(this.Out)
			templates.NewController(this.Package, this.Name, this.Out).Generate()
		} else if enum == this.Gen {
			if this.Name == "" {
				log.Fatalln(`enum configure file unable empty 
eg: -name=test\bin\enum_cmd.md 
or format
eg: -name="-e=state -f=状态:issue-1-发布,draft-2-草稿"`)
			}
			if this.Out == "" {
				log.Fatalln(`output directory unable empty 
eg: -out=test\app\Enums`)
			}
			this.Package = utils.Basename(this.Out)
			templates.NewEnum(this.Package, this.Name, this.Out)//.Generate()
		} else if model == this.Gen {
			if this.Name == "" {
				log.Fatalln(`table name unable empty 
eg: -name=ts_user`)
			}
			if this.Out == "" {
				log.Fatalln(`output directory unable empty 
eg: -out=test\app\Models`)
			}
			db := mapperOrm().DB
			if this.Name == "all" {
				newModel := templates.NewModel(db, this.Name, this.Out, GetDbConfig().Database, GetDbConfig().Prefix)
				tables := newModel.GetTables()
				for _, table := range tables {
					templates.NewModel(db, table.Name, this.Out, GetDbConfig().Database, GetDbConfig().Prefix).Generate()
				}
			} else {
				templates.NewModel(db, this.Name, this.Out, GetDbConfig().Database, GetDbConfig().Prefix).Generate()
			}
		} else {
			log.Fatalln(`gen Arguments Error! 
Explain: Generate Controller or Model 
--option[controller | model | enum] 
eg:-gen=controller`)
		}
		os.Exit(1)
	}
}
