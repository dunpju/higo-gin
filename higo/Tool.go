package higo

import (
	"flag"
	"fmt"
	"github.com/dengpju/higo-gin/higo/templates"
	"github.com/dengpju/higo-utils/utils"
	"log"
	"os"
	"regexp"
)

const (
	controller = "controller"
	model      = "model"
	enum       = "enum"
	codes      = "code"
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
		flag.StringVar(&this.Gen, "gen", "", `explain: Generate Controller or Model or Enum or Code or Dao or Entity
	--option[controller | model | enum | code | dao | entity]
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
			templates.NewEnum(this.Package, this.Name, this.Out).Generate()
		} else if codes == this.Gen {
			if this.Name == "" {
				log.Fatalln(`code configure file unable empty 
	eg: -name=test\bin\code_cmd.md 
	or format
	eg: -name="-c=token -i=400001 -f=token码:token_empty-token为空"`)
			}
			if this.Out == "" {
				log.Fatalln(`output directory unable empty 
	eg: -out=test\app\Codes`)
			}
			this.Package = utils.Basename(this.Out)
			templates.NewCode(this.Package, this.Name, this.Out).Generate()
		} else if model == this.Gen {
			if this.Name == "" {
				log.Fatalln(`table name unable empty 
eg: -name=ts_user`)
			}
			if this.Out == "" {
				log.Fatalln(`output directory unable empty 
eg: -out=test\app\Models`)
			}
		loopDao:
			fmt.Print("Whether To Generate Dao [yes|no] (default:yes):")
			confirmBeginGenerate := "yes"
			isGenerateDao := "yes"
			isGenerateEntity := "yes"
			n, err := fmt.Scanln(&isGenerateDao)
			if nil != err && n > 0 {
				panic(err)
			}
			if ("yes" != isGenerateDao && "no" != isGenerateDao) && n > 0 {
				goto loopDao
			}
			fmt.Printf("Your Choice Generate Dao: %s\n", isGenerateDao)
			capitalBeganReg := regexp.MustCompile(`^[A-Z].*`)
			if capitalBeganReg == nil {
				log.Fatalln("regexp err")
			}
			daoDir := "dao"
			isMatchCapitalBegan := capitalBeganReg.FindString(utils.Basename(this.Out))
			if isMatchCapitalBegan != "" {
				daoDir = utils.Ucfirst(daoDir)
			}
			daoDir = utils.Dirname(this.Out) + `\` + daoDir
			fmt.Printf("Whether To Confirm Output Directory Of Dao Default (%s)? Enter/Input: ", daoDir)
			inputDaoDir := daoDir
			n, err = fmt.Scanln(&inputDaoDir)
			if nil != err && n > 0 {
				panic(err)
			}
			fmt.Printf("You Confirmed Output Directory Of Dao: %s\n", inputDaoDir)
		loopEntity:
			fmt.Print("Whether To Generate Entity [yes|no] (default:yes):")
			n, err = fmt.Scanln(&isGenerateEntity)
			if nil != err && n > 0 {
				panic(err)
			}
			if ("yes" != isGenerateEntity && "no" != isGenerateEntity) && n > 0 {
				goto loopEntity
			}
			fmt.Printf("Your Choice Generate Entity: %s\n", isGenerateEntity)
			entityDir := "entity"
			isMatchCapitalBegan = capitalBeganReg.FindString(utils.Basename(this.Out))
			if isMatchCapitalBegan != "" {
				entityDir = utils.Ucfirst(entityDir)
			}
			entityDir = utils.Dirname(this.Out) + `\` + entityDir
			fmt.Printf("Whether To Confirm Output Directory Of Entity Default (%s)? Enter/Input: ", entityDir)
			inputEntityDir := entityDir
			n, err = fmt.Scanln(&inputEntityDir)
			if nil != err && n > 0 {
				panic(err)
			}
			fmt.Printf("You Confirmed Output Directory Of Entity: %s\n", inputEntityDir)
			//确认开始构建
		loopConfirmBeginGenerate:
			fmt.Print("Confirm To Start Generate [yes|no] (default:yes):")
			n, err = fmt.Scanln(&confirmBeginGenerate)
			if ("yes" != confirmBeginGenerate && "no" != confirmBeginGenerate) && n > 0 {
				goto loopConfirmBeginGenerate
			}
			if ("yes" != confirmBeginGenerate) && n > 0 {
				goto loopDao
			}
			fmt.Print("Start Generate ......")
			//连接数据库准备构建
			db := newOrm().DB
			if this.Name == "all" {
				tables := templates.GetTables(db, GetDbConfig().Database)
				for _, table := range tables {
					genModel := templates.NewModel(db, table.Name, this.Out, GetDbConfig().Database, GetDbConfig().Prefix)
					genModel.Generate()
				}
			} else {
				genModel := templates.NewModel(db, this.Name, this.Out, GetDbConfig().Database, GetDbConfig().Prefix)
				genModel.Generate()
				log.Fatalln(genModel)
			}
		} else {
			log.Fatalln(`gen Arguments Error! 
Explain: Generate Controller or Model or Enum or Code or Dao or Entity
	--option[controller | model | enum | code | dao | entity] 
	eg:-gen=controller`)
		}
		os.Exit(1)
	}
}
