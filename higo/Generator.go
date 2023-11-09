package higo

import (
	"fmt"
	"github.com/dunpju/higo-gin/higo/templates"
	"github.com/dunpju/higo-orm/gen"
	"github.com/dunpju/higo-utils/utils"
	"github.com/dunpju/higo-utils/utils/dirutil"
	"github.com/dunpju/higo-utils/utils/stringutil"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	gen.InitModel()
	InitController()
	genCommand.AddCommand(gen.ModelCommand)
	genCommand.AddCommand(ControllerCommand)
	rootCommand.AddCommand(genCommand)
}

var (
	name string
	out  string
)

var rootCommand = &cobra.Command{
	Use:   "",
	Short: "命令行工具",
	Long:  `命令行工具`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("命令行工具")
	},
}

var genCommand = &cobra.Command{
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
	Example: "controller --name=Test --out=app\\controllers",
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
