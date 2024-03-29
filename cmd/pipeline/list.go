package pipeline

import (
	"bb/api"
	"bb/util"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pipelines from a repository",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		nResults, _ := cmd.Flags().GetInt("number-results")
		showAuthor, _ := cmd.Flags().GetBool("author")
		targetBranch, _ := cmd.Flags().GetString("target")

		for pipeline := range api.GetPipelineList(viper.GetString("repo"), nResults, targetBranch) {
			if pipeline.State.Result.Name == "" {
				util.Printf("%s", util.FormatPipelineStatus(pipeline.State.Name))
			} else {
				util.Printf("%s", util.FormatPipelineStatus(pipeline.State.Result.Name))
			}
			util.Printf(" \033[1;32m#%d\033[m ", pipeline.BuildNumber)
			if pipeline.Target.Source != "" {
				util.Printf("%s \033[1;34m[ %s → %s ]\033[m", pipeline.Target.PullRequest.Title, pipeline.Target.Source, pipeline.Target.Destination)
			} else {
				util.Printf("\033[1;34m[ %s ]\033[m", pipeline.Target.RefName)
			}

			util.Printf(" \033[37m%s (%s)\033[m\n", util.TimeDuration(time.Duration(pipeline.DurationInSeconds*1000000000)), util.TimeAgo(pipeline.CreatedOn))

			if showAuthor {
				util.Printf("        \033[33m%s\033[m \033[37mTrigger: %s\033[m\n", pipeline.Author.DisplayName, pipeline.Trigger.Name) //  \033[37mComments: %d\033[m",
			}
		}
	},
}

func init() {
	ListCmd.Flags().IntP("number-results", "n", 10, "max number of results retrieve (max: 100)")
	ListCmd.Flags().BoolP("author", "a", false, "show author information")
	ListCmd.Flags().String("target", "", "filter target branch of pipeline")
	ListCmd.RegisterFlagCompletionFunc("target", util.BranchCompletion)
}
