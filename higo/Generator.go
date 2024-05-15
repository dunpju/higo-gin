package higo

import (
	"fmt"
	gen2 "github.com/dunpju/higo-enum/gen"
	"github.com/dunpju/higo-gin/higo/templates"
	"github.com/dunpju/higo-orm/gen"
	"github.com/dunpju/higo-utils/utils"
	"github.com/dunpju/higo-utils/utils/dirutil"
	"github.com/dunpju/higo-utils/utils/stringutil"
	"github.com/spf13/cobra"
	"log"
	"os"
	"regexp"
)

var (
	name            string
	out             string
	yamlPath        string
	jsonFile        string
	tag             string
	isYaml          bool
	isAuto          bool
	isForce         bool
	capitalBeganReg = regexp.MustCompile(`^[A-Z].*`) //匹配大写字母开头
)

func RootInit() {
	GenCommandInit()
	RootCommand.AddCommand(GenCommand)
}

var RootCommand = &cobra.Command{
	Use:   "",
	Short: "命令行工具",
	Long:  `命令行工具`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("命令行工具")
	},
}

func GenCommandInit() {
	gen.InitModel()
	InitController()
	InitParam()
	InitCode()
	gen2.InitEnum()
	InitService()
	GenCommand.AddCommand(gen.ModelCommand)
	GenCommand.AddCommand(ControllerCommand)
	GenCommand.AddCommand(ParamCommand)
	GenCommand.AddCommand(CodeCommand)
	GenCommand.AddCommand(gen2.EnumCommand)
	GenCommand.AddCommand(ServiceCommand)
}

var GenCommand = &cobra.Command{
	Use:     "gen",
	Short:   "构建工具",
	Long:    "构建工具",
	Example: "gen",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func InitController() {
	ControllerCommand.Flags().StringVarP(&name, "name", "n", "", "控制器名称")
	err := ControllerCommand.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}
	ControllerCommand.Flags().StringVarP(&out, "out", "o", "", "生成目录,如:app\\controllers")
	err = ControllerCommand.MarkFlagRequired("out")
	if err != nil {
		panic(err)
	}
}

var ControllerCommand = &cobra.Command{
	Use:     "controller",
	Short:   "Controller构建工具",
	Long:    "Controller构建工具",
	Example: "controller --name=Test --out=app\\Controllers",
	Run: func(cmd *cobra.Command, args []string) {
	loopParam:
		isMatchCapitalBegan := ""
		controllerTool := templates.NewControllerTool()
		controllerTool.Name = name
		controllerTool.OutParamDir = name
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
			outParamDir := utils.Dir.Dirname(out) + `\` + "params"
			fmt.Printf("Confirm Output Directory Of Param Default (%s)? Enter/Input: ", outParamDir)
			controllerTool.OutParamDir = outParamDir
			n, err = fmt.Scanln(&controllerTool.OutParamDir)
			if nil != err && n > 0 {
				panic(err)
			}
			fmt.Printf("Confirmed Output Directory Of Param: %s\n", controllerTool.OutParamDir)
		}
		fmt.Print("Start Generate ......\n")
		pkg := dirutil.Basename(out)
		templates.NewController(pkg, name, out).Generate()
		controllerTool.Generate()
	},
}

func InitParam() {
	ParamCommand.Flags().StringVarP(&name, "name", "n", "", "参数结构体名称")
	err := ParamCommand.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}
	ParamCommand.Flags().StringVarP(&out, "out", "o", "", "生成目录,如:app\\params")
	err = ParamCommand.MarkFlagRequired("out")
	if err != nil {
		panic(err)
	}
	ParamCommand.Flags().StringVarP(&jsonFile, "json", "j", "", "输入json文件")
	ParamCommand.Flags().BoolVarP(&isForce, "force", "f", false, "强制覆盖")
	ParamCommand.Flags().StringVarP(&tag, "tag", "t", "", "tag,如json")
}

var ParamCommand = &cobra.Command{
	Use:   "param",
	Short: "Param构建工具",
	Long:  "Param构建工具",
	Example: `param --name=list --out=app\Params
param --name=list --out=app\Params --json=bin\param.json --force --tag=json`,
	Run: func(cmd *cobra.Command, args []string) {
		templates.NewParam(name, out, jsonFile, tag, isForce).Generate()
	},
}

func InitCode() {
	CodeCommand.Flags().StringVarP(&name, "name", "n", "", "Code结构体名称,如:ErrorCode")
	err := CodeCommand.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}
	CodeCommand.Flags().StringVarP(&out, "out", "o", "", "生成目录,如:app\\Codes")
	err = CodeCommand.MarkFlagRequired("out")
	if err != nil {
		panic(err)
	}
	CodeCommand.Flags().StringVarP(&yamlPath, "path", "p", "", "yaml配置文件路径,如:bin\\200.yaml或bin\\yaml目录")
	err = CodeCommand.MarkFlagRequired("out")
	if err != nil {
		panic(err)
	}
	CodeCommand.Flags().BoolVarP(&isYaml, "yaml", "y", false, "生成yaml模板文件")
	CodeCommand.Flags().BoolVarP(&isAuto, "auto", "a", false, "自动增长")
	CodeCommand.Flags().BoolVarP(&isForce, "force", "f", false, "强制覆盖")
}

var CodeCommand = &cobra.Command{
	Use:   "code",
	Short: "Code构建工具",
	Long:  "Code构建工具",
	Example: `  - 1、单配置文件生成: 
    code --name=ErrorCode --out=app\Codes --path=bin\200.yaml
  - 2、多配置文件,覆盖生成: 
    code --name=ErrorCode --out=app\Codes --path=bin\yaml --auto --force
  - 3、生成yaml模版文件: 
    code --name=ErrorCode --out=app\Codes --path=bin\200.yaml --yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		if isYaml {
			matchYamlSuffix, err := regexp.Match(`\.yaml`, []byte(yamlPath))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if !matchYamlSuffix {
				yamlPath += ".yaml"
			}
			yamlTpl := utils.Dir.Basename(yamlPath, ".yaml")
			match, err := regexp.Match(`\d+`, []byte(yamlTpl))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if !match {
				fmt.Println(fmt.Errorf("生成yaml文件,path必须指定数字文件名,如200.yaml"))
				os.Exit(1)
			}
			if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
				utils.Dir.Mkdir(yamlPath, os.ModePerm)
			}
			f, err := os.Create(yamlPath)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			fileContext := "tpl_code_模板code请修改:\n"
			fileContext += "  message: \"模板消息请修改\"\n"
			fileContext += "  iota: \"yes\"\n"
			_, err = f.Write([]byte(fileContext))
			if err != nil {
				panic(err)
			}
		} else {
			codeArguments := &templates.CodeArguments{
				Package: dirutil.Basename(out),
				Name:    name,
				Out:     out,
				Path:    yamlPath,
			}
			codeArguments.Auto = gen.No
			if isAuto {
				codeArguments.Auto = gen.Yes
			}
			codeArguments.Force = gen.No
			if isForce {
				codeArguments.Force = gen.Yes
			}
			templates.NewCode(codeArguments).Generate()
		}
	},
}

func InitService() {
	ServiceCommand.Flags().StringVarP(&name, "name", "n", "", `服务名称,--name=TestService`)
	err := ServiceCommand.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}
	ServiceCommand.Flags().StringVarP(&out, "out", "o", "", "生成目录,如:app\\Services")
	err = ServiceCommand.MarkFlagRequired("out")
	if err != nil {
		panic(err)
	}
}

var ServiceCommand = &cobra.Command{
	Use:     "service",
	Short:   "Service构建工具",
	Long:    "Service构建工具",
	Example: "service --name=TestService --out=app\\Services",
	Run: func(cmd *cobra.Command, args []string) {
		templates.NewService(dirutil.Basename(out), name, out).Generate()
	},
}
