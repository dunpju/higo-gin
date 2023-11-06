package higo

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Generator = &cobra.Command{
	Use:   "",
	Short: "构建工具",
	Long:  `构建工具`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("build tools")
	},
}
