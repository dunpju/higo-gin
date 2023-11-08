package higo

import (
	"fmt"
	"github.com/dunpju/higo-orm/gen"
	"github.com/spf13/cobra"
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

func init() {
	gen.InitModel()
	genCommand.AddCommand(gen.ModelCommand)
	rootCommand.AddCommand(genCommand)
}
