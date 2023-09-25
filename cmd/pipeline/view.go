package pipeline

import (
	"bb/api"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ViewCmd = &cobra.Command{
	Use:   "view [ID]",
	Short: "View details of a pipeline",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api.GetPipelineList(viper.GetString("repo"), 0)
		fmt.Println("Not implemented")
	},
}

func init() {
}
