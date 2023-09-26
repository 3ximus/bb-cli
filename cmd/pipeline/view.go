package pipeline

import (
	"bb/api"
	"bb/util"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ViewCmd = &cobra.Command{
	Use:   "view [ID]",
	Short: "View details of a pipeline",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")

		var id int
		var err error
		if len(args) == 0 {
			branch, err := util.GetCurrentBranch()
			cobra.CheckErr(err)
			// retrieve id of pr for current branch
			pipeline := <-api.GetPipelineList(repo, 1, branch)
			if pipeline.BuildNumber == 0 {
				cobra.CheckErr("No pipelines found for this branch")
			}
			id = pipeline.BuildNumber
		} else {
			id, err = strconv.Atoi(args[0])
			cobra.CheckErr(err)
		}

		pipeline := <-api.GetPipeline(repo, id)

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

		fmt.Printf("        \033[33m%s\033[m \033[37mTrigger: %s\033[m\n", pipeline.Author.DisplayName, pipeline.Trigger.Name) //  \033[37mComments: %d\033[m",

		web, _ := cmd.Flags().GetBool("web")
		if web {
			util.OpenInBrowser(api.BBBrowsePipelines(repo, id))
			return
		}

	},
}

func init() {
	ViewCmd.Flags().Bool("web", false, "Open in the browser.")
}
