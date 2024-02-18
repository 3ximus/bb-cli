package pipeline

import (
	"bb/api"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var VariablesCmd = &cobra.Command{
	Use:     "variables",
	Short:   "List pipeline variables",
	Long:    "List pipeline variables. If variable is secured only *** is displayed",
	Aliases: []string{"var"},
	Run: func(cmd *cobra.Command, args []string) {
		for variable := range api.GetPipelineVariables(viper.GetString("repo")) {
			if variable.Secured {
				fmt.Printf("%s = \033[37m***\033[m", variable.Key)
			} else {
				fmt.Printf("%s = \033[37m%s\033[m", variable.Key, variable.Value)
			}
			fmt.Println()
		}
	},
}
