package higo

import (
	"flag"
	"fmt"
	"github.com/dunpju/higo-gin/higo/templates"
	"github.com/dunpju/higo-orm/gen"
	"github.com/dunpju/higo-utils/utils"
	"github.com/dunpju/higo-utils/utils/dirutil"
	"github.com/dunpju/higo-utils/utils/stringutil"
	"log"
	"os"
	"regexp"
)

const (
	service    = "service"
	controller = "controller"
	model      = "model"
	enum       = "enum"
	codes      = "code"
	param      = "param"
)

var capitalBeganReg = regexp.MustCompile(`^[A-Z].*`) //匹配大写字母开头

type Tool struct {
	Gen     string
	Name    string
	Out     string
	Package string
	Auto    string
	Force   string
	Path    string
	Const   string
	Code    string
	Message string
	Iota    string
	Connect string
}

func NewTool() *Tool {
	return &Tool{}
}

func (this *Tool) Parse() *ToolParse {
	flag.StringVar(&this.Gen, "gen", "nil", `explain: Generate Controller Model Enum Code Dao Entity
    --option[controller | model | enum | code | dao | entity | param | service]
    eg:-gen=controller`)
	flag.StringVar(&this.Name, "name", "", `explain: Generate Name 
    eg:-name=Test`)
	flag.StringVar(&this.Out, "out", "", `explain: Generate file output directory 
    eg:-out=app\Controllers`)
	flag.StringVar(&this.Auto, "auto", "no", `explain: Code Autoloading yes/no`)
	flag.StringVar(&this.Force, "force", "no", `explain: Forced updating yes/no`)
	flag.StringVar(&this.Path, "path", "", `explain: A configuration file of code path`)
	flag.StringVar(&this.Const, "const", "", `explain: Code const`)
	flag.StringVar(&this.Code, "code", "", `explain: Code number`)
	flag.StringVar(&this.Message, "message", "", `explain: Code message`)
	flag.StringVar(&this.Iota, "iota", "no", `explain: Code iota yes/no`)
	flag.StringVar(&this.Connect, "connect", "", `explain: Default`)
	flag.Parse()
	return newToolParse(this)
}

type ToolParse struct {
	tool *Tool
}

func newToolParse(tool *Tool) *ToolParse {
	return &ToolParse{tool: tool}
}

func (t *ToolParse) Cmd() {
	t.tool.Cmd()
}

func (this *Tool) Cmd() {
	if this.Gen != "" && this.Gen != "nil" {
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
		} else if service == this.Gen {
			this.service()
		} else {
			log.Fatalln(`gen Arguments Error! 
Explain: Generate Controller/Model/Enum/Code/Dao/Entity/Param/Service
    --option[controller | model | enum | code | dao | entity | param | service]
    eg:-gen=controller
    Tips*: Dao and Entity with Model, It's a composite generate tool; 
        They using AST parsing source files, the original code will not be affected; 
        So perform model generate`)
		}
		os.Exit(1)
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
loopParam:
	isMatchCapitalBegan := ""
	controllerTool := templates.NewControllerTool()
	controllerTool.Name = this.Name
	controllerTool.OutParamDir = this.Name
	fmt.Print("Generate Param List/Add/Edit/Delete [yes|no] (default:yes):")
	n, err := fmt.Scanln(&controllerTool.IsGenerateParam)
	if nil != err && n > 0 {
		panic(err)
	}
	if (gen.Yes != controllerTool.IsGenerateParam && gen.No != controllerTool.IsGenerateParam) && n > 0 {
		goto loopParam
	}
	fmt.Printf("Choice Generate Param: %s\n", controllerTool.IsGenerateParam)
	if controllerTool.IsGenerateParam.Bool() { // 确认构建param
		controllerTool.ParamTag = append(controllerTool.ParamTag, templates.List, templates.Add, templates.Edit, templates.Delete)
		if capitalBeganReg == nil {
			log.Fatalln("regexp err")
		}
		isMatchCapitalBegan = capitalBeganReg.FindString(controllerTool.Name)
		if isMatchCapitalBegan != "" {
			controllerTool.Name = stringutil.Ucfirst(controllerTool.Name)
		}
		outParamDir := utils.Dir.Dirname(this.Out) + `\` + "params"
		fmt.Printf("Confirm Output Directory Of Param Default (%s)? Enter/Input: ", outParamDir)
		controllerTool.OutParamDir = outParamDir
		n, err = fmt.Scanln(&controllerTool.OutParamDir)
		if nil != err && n > 0 {
			panic(err)
		}
		fmt.Printf("Confirmed Output Directory Of Param: %s\n", controllerTool.OutParamDir)
	}
	fmt.Print("Start Generate ......\n")
	this.Package = dirutil.Basename(this.Out)
	templates.NewController(this.Package, this.Name, this.Out).Generate()
	controllerTool.Generate()
}

func (this *Tool) service() {
	if this.Name == "" {
		log.Fatalln(`service name unable empty 
    eg: -name=TestService`)
	}
	if this.Out == "" {
		log.Fatalln(`output directory unable empty 
    eg: -out=app\Services`)
	}
	this.Package = dirutil.Basename(this.Out)
	templates.NewService(this.Package, this.Name, this.Out).Generate()
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
		log.Fatalln(`code struct name unable empty
    eg: -name=ErrorCode`)
	}
	if this.Out == "" {
		log.Fatalln(`output directory unable empty 
    eg: -out=app\Codes`)
	}
	if this.Path == "" {
		if this.Const == "" && this.Code == "" && this.Message == "" {
			log.Fatalln(`a configuration file of code path unable empty or code const unable empty
    eg: -path=bin\200.yaml a file Or bin\yaml a directory of yaml file
    yaml file format
success:
  code: 200
  message: "成功"`)
		} else if this.Const == "" {
			log.Fatalln(`code const unable empty 
    eg: -const=success`)
		} else if this.Code == "" {
			log.Fatalln(`code number unable empty 
    eg: -code=200`)
		} else if this.Message == "" {
			log.Fatalln(`code message unable empty 
    eg: -message=成功`)
		}
	}
	this.Package = dirutil.Basename(this.Out)
	codeArguments := &templates.CodeArguments{
		Package: this.Package,
		Name:    this.Name,
		Out:     this.Out,
		Auto:    this.Auto,
		Force:   this.Force,
		Path:    this.Path,
		Const:   this.Const,
		Code:    this.Code,
		Message: this.Message,
		Iota:    this.Iota,
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
	if this.Connect == "" {
		log.Fatalln(`data source connect name unable empty 
    eg: -connect=Default`)
	}
	if this.Name == "" {
		log.Fatalln(`table name unable empty 
    eg: -name=ts_user`)
	}
	if this.Out == "" {
		log.Fatalln(`output directory unable empty 
    eg: -out=app\Models`)
	}
}
