package pipeline

import (
	"bb/api"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RunCmd = &cobra.Command{
	Use:   "run [BRANCH]",
	Short: "Run pipeline for branch",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api.GetPipelineList(viper.GetString("repo"), 0)
		fmt.Println("Not implemented")
	},
}

func init() {
}
