package pipeline

import (
	"bb/api"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running pipeline",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		api.GetPipelineList(viper.GetString("repo"), 0, "")
		fmt.Println("Not implemented")
	},
}

func init() {
}
