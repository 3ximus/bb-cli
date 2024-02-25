package environment

import (
	"bb/api"
	"bb/util"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var VariablesCmd = &cobra.Command{
	Use:     "variables",
	Short:   "List variables for specific environment",
	Long:    "List variables for specific environment. If variable is secured only *** is displayed",
	Aliases: []string{"var"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for variable := range api.GetEnvironmentVariables(viper.GetString("repo"), args[0]) {
			if variable.Secured {
				util.Printf("%s = \033[37m***\033[m", variable.Key)
			} else {
				util.Printf("%s = \033[37m%s\033[m", variable.Key, variable.Value)
			}

			fmt.Println()
		}
	},
}
