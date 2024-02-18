package pipeline

import (
	"bb/api"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var LogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show logs of a pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		api.GetPipelineList(viper.GetString("repo"), 0, "")
		fmt.Println("Not implemented")
	},
}

func init() {
}
