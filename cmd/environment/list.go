package environment

import (
	"bb/api"
	"bb/util"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List environments from a repository",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		status, _ := cmd.Flags().GetBool("status")
		for environment := range api.GetEnvironmentList(viper.GetString("repo"), status) {
			if status {
				if environment.Status.State.Result.Name == "" {
					util.Printf("%s ", util.FormatPipelineStatus(environment.Status.State.Name))
				} else {
					util.Printf("%s ", util.FormatPipelineStatus(environment.Status.State.Result.Name))
				}
			}
			util.Printf("%s \033[37m%s\033[m", environment.Name, environment.EnvironmentType.Name)
			fmt.Println()
		}
	},
}

func init() {
	ListCmd.Flags().BoolP("status", "s", false, "show deployment status of each environment")
}
