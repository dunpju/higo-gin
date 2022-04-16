package higo

import (
	"flag"
	"fmt"
	"github.com/dengpju/higo-gin/higo/templates"
	"github.com/dengpju/higo-utils/utils/dirutil"
	"github.com/dengpju/higo-utils/utils/stringutil"
	"log"
	"os"
	"regexp"
)

const (
	controller                 = "controller"
	model                      = "model"
	enum                       = "enum"
	codes                      = "code"
	param                      = "param"
	yes        templates.YesNo = "yes"
	no         templates.YesNo = "no"
)

type Tool struct {
	Gen     string
	Name    string
	Out     string
	Package string
	Auto    string
	Force   string
	Extends string
}

func NewTool() *Tool {
	return &Tool{}
}

func (this *Tool) Cmd() {
	if len(os.Args) >= 2 {
		flag.StringVar(&this.Gen, "gen", "", `explain: Generate Controller Model Enum Code Dao Entity
    --option[controller | model | enum | code | dao | entity | param]
    eg:-gen=controller`)
		flag.StringVar(&this.Name, "name", "", `explain: Generate Name 
    eg:-name=Test`)
		flag.StringVar(&this.Out, "out", "", `explain: Generate file output directory 
    eg:-out=app\Controllers`)
		flag.StringVar(&this.Auto, "auto", "no", `explain: Autoloading`)
		flag.StringVar(&this.Force, "force", "no", `explain: Forced updating`)
		flag.StringVar(&this.Extends, "extends", "", `explain: Extends struct`)
		flag.Parse()
		if controller == this.Gen {
			this.controller()
		} else if enum == this.Gen {
			this.enum()
		} else if codes == this.Gen {
			this.code()
		} else if param == this.Gen {
			this.param()
		} else if model == this.Gen {
			this.model()
		} else {
			log.Fatalln(`gen Arguments Error! 
Explain: Generate Controller Model Enum Code Dao Entity
    --option[controller | model | enum | code | dao | entity | param]
    eg:-gen=controller`)
		}
		os.Exit(1)
	}
}

func daoEntity(modelTool templates.ModelTool, genModel templates.Model) {
	if modelTool.IsGenerateDao.Bool() {
		entity := templates.NewEntity(modelTool, genModel)
		entity.Generate()
		templates.NewDao(modelTool, genModel, *entity).Generate()
	} else if modelTool.IsGenerateEntity.Bool() {
		templates.NewEntity(modelTool, genModel).Generate()
	}
}

func (this *Tool) controller() {
	if this.Name == "" {
		log.Fatalln(`controller name unable empty 
    eg: -name=Test`)
	}
	if this.Out == "" {
		log.Fatalln(`output directory unable empty 
    eg: -out=app\Controllers`)
	}
	this.Package = dirutil.Basename(this.Out)
	templates.NewController(this.Package, this.Name, this.Out).Generate()
}

func (this *Tool) enum() {
	if this.Name == "" {
		log.Fatalln(`enum configure file unable empty 
    eg: -name=bin\enum_cmd.md
    eg: -name="-e=state -f=状态:issue-1-发布,draft-2-草稿"`)
	}
	if this.Out == "" {
		log.Fatalln(`output directory unable empty 
    eg: -out=app\Enums`)
	}
	templates.NewEnum(this.Name, this.Out).Generate()
}

func (this *Tool) code() {
	if this.Name == "" {
		log.Fatalln(`code configure file unable empty 
    eg: -name="bin\code_cmd.md"
    eg: -name="bin\200.yaml"
    eg: -name="bin\yaml" a directory
    eg: -name="-c=token -i=400001 -f=token码:token_empty-token为空"`)
	}
	if this.Out == "" {
		log.Fatalln(`output directory unable empty 
    eg: -out=app\Codes`)
	}
	if this.Extends == "" {
		log.Fatalln(`extends struct unable empty 
    eg: -extends=CodeErrorCode`)
	}
	this.Package = dirutil.Basename(this.Out)
	codeArguments := &templates.CodeArguments{
		Package: this.Package,
		Name:    this.Name,
		Out:     this.Out,
		Auto:    this.Auto,
		Force:   this.Force,
		Extends: this.Extends,
	}
	templates.NewCode(codeArguments).Generate()
}

func (this *Tool) param() {
	if this.Name == "" {
		log.Fatalln(`param name unable empty 
    eg: -name=list`)
	}
	if this.Out == "" {
		log.Fatalln(`output directory unable empty 
    eg: -out=app\Params`)
	}
	templates.NewParam(this.Name, this.Out).Generate()
}

func (this *Tool) model() {
	if this.Name == "" {
		log.Fatalln(`table name unable empty 
    eg: -name=ts_user`)
	}
	if this.Out == "" {
		log.Fatalln(`output directory unable empty 
    eg: -out=app\Models`)
	}
loopDao:
	capitalBeganReg := regexp.MustCompile(`^[A-Z].*`) //匹配大写字母开头
	isMatchCapitalBegan := ""
	modelTool := templates.NewModelTool()
	modelTool.Name = this.Name
	modelTool.Out = this.Out
	fmt.Print("Generate Dao [yes|no] (default:yes):")
	n, err := fmt.Scanln(&modelTool.IsGenerateDao)
	if nil != err && n > 0 {
		panic(err)
	}
	if (yes != modelTool.IsGenerateDao && no != modelTool.IsGenerateDao) && n > 0 {
		goto loopDao
	}
	fmt.Printf("Choice Generate Dao: %s\n", modelTool.IsGenerateDao)
	if modelTool.IsGenerateDao.Bool() { //确认构建dao
		if capitalBeganReg == nil {
			log.Fatalln("regexp err")
		}
		daoDir := "dao"
		isMatchCapitalBegan = capitalBeganReg.FindString(dirutil.Basename(this.Out))
		if isMatchCapitalBegan != "" {
			daoDir = stringutil.Ucfirst(daoDir)
		}
		outDaoDir := dirutil.Dirname(this.Out) + `\` + daoDir
		fmt.Printf("Confirm Output Directory Of Dao Default (%s)? Enter/Input: ", outDaoDir)
		modelTool.OutDaoDir = outDaoDir
		n, err = fmt.Scanln(&modelTool.OutDaoDir)
		if nil != err && n > 0 {
			panic(err)
		}
		fmt.Printf("Confirmed Output Directory Of Dao: %s\n", modelTool.OutDaoDir)
		//确认构建dao，默认必须构建entity
		modelTool.IsGenerateEntity = yes
		goto loopChoiceGenerateEntity
	}
loopEntity:
	fmt.Print("Generate Entity [yes|no] (default:yes):")
	n, err = fmt.Scanln(&modelTool.IsGenerateEntity)
	if nil != err && n > 0 {
		panic(err)
	}
	if (yes != modelTool.IsGenerateEntity && no != modelTool.IsGenerateEntity) && n > 0 {
		goto loopEntity
	}
loopChoiceGenerateEntity:
	fmt.Printf("Choice Generate Entity: %s\n", modelTool.IsGenerateEntity)
	if modelTool.IsGenerateEntity.Bool() { //确认构建entity
		entityDir := "entity"
		isMatchCapitalBegan = capitalBeganReg.FindString(dirutil.Basename(this.Out))
		if isMatchCapitalBegan != "" {
			entityDir = stringutil.Ucfirst(entityDir)
		}
		outEntityDir := dirutil.Dirname(this.Out) + `\` + entityDir
		fmt.Printf("Confirm Output Directory Of Entity Default (%s)? Enter/Input: ", outEntityDir)
		modelTool.OutEntityDir = outEntityDir
		n, err = fmt.Scanln(&modelTool.OutEntityDir)
		if nil != err && n > 0 {
			panic(err)
		}
		fmt.Printf("Confirmed Output Directory Of Entity: %s\n", modelTool.OutEntityDir)
	}
	//确认开始构建
loopConfirmBeginGenerate:
	fmt.Print("Start Generate [yes|no] (default:yes):")
	n, err = fmt.Scanln(&modelTool.ConfirmBeginGenerate)
	if (yes != modelTool.ConfirmBeginGenerate && no != modelTool.ConfirmBeginGenerate) && n > 0 {
		goto loopConfirmBeginGenerate
	}
	if (yes != modelTool.ConfirmBeginGenerate) && n > 0 {
		goto loopDao
	}
	fmt.Println(modelTool)
	fmt.Print("Start Generate ......")
	//连接数据库准备构建
	db := newOrm().DB
	if this.Name == "all" {
		tables := templates.GetDbTables(db, GetDbConfig().Database)
		for _, table := range tables {
			genModel := templates.NewModel(db, table.Name, this.Out, GetDbConfig().Database, GetDbConfig().Prefix)
			genModel.Generate()
			daoEntity(*modelTool, *genModel)
			fmt.Println("==================================================================================")
		}
	} else {
		genModel := templates.NewModel(db, this.Name, this.Out, GetDbConfig().Database, GetDbConfig().Prefix)
		genModel.Generate()
		daoEntity(*modelTool, *genModel)
	}
}
