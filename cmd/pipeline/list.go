package pipeline

import (
	"bb/api"
	"bb/util"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pipelines from a repository",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		pages, _ := cmd.Flags().GetInt("pages")
		pipelineChannel := api.GetPipelineList(viper.GetString("repo"), pages)
		for pipeline := range pipelineChannel {
			if pipeline.State.Result.Name == "" {
				fmt.Printf(" %s", util.FormatPipelineState(pipeline.State.Name))
			} else {
				fmt.Printf(" %s", util.FormatPipelineState(pipeline.State.Result.Name))
			}
			fmt.Printf(" \033[1;32m#%d\033[m ", pipeline.BuildNumber)
			if pipeline.Target.Source != "" {
				fmt.Printf("%s \033[1;34m[ %s → %s]\033[m\n", pipeline.Target.PullRequest.Title, pipeline.Target.Source, pipeline.Target.Destination)
			} else {
				fmt.Printf("\033[1;34m[ %s ]\033[m\n", pipeline.Target.RefName)
			}

			fmt.Printf("       \033[33m%s\033[m\n", pipeline.Author.DisplayName) //  \033[37mComments: %d\033[m",
		}
	},
}

func init() {
	ListCmd.Flags().Int("pages", 1, "number of pages with results to retrieve")
}
