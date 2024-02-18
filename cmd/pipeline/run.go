package pipeline

import (
	"bb/api"
	"bb/util"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RunCmd = &cobra.Command{
	Use:   "run [BRANCH]",
	Short: "Run pipeline for branch",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")
		if branch == "" {
			if len(args) == 0 {
				var err error
				branch, err = util.GetCurrentBranch()
				cobra.CheckErr(err)
			} else {
				branch = args[0]
			}
		}

		newpipeline := api.RunPipelineRequestBody{}
		newpipeline.Target.RefName = branch
		newpipeline.Target.Type = "pipeline_ref_target"
		newpipeline.Target.RefType = "branch"

		selectedConfig, _ := cmd.Flags().GetString("select")
		if selectedConfig != "" {
			newpipeline.Target.Selector = &api.PipelineSelectorBody{}
			newpipeline.Target.Selector.Type = "custom"
			newpipeline.Target.Selector.Pattern = selectedConfig
		}

		commit, _ := cmd.Flags().GetString("commit")
		if commit != "" {
			newpipeline.Target.Commit = &api.PipelineCommitRefBody{}
			newpipeline.Target.Commit.Hash = commit
			newpipeline.Target.Commit.Type = "commit"
		}

		// TODO untested
		pullRequest, _ := cmd.Flags().GetString("pull-request")
		if pullRequest != "" {
			newpipeline.Target.PullRequest = &api.PipelinePullRequestBody{}
			newpipeline.Target.PullRequest.Id = pullRequest
			newpipeline.Target.Type = "pipeline_pullrequest_target"
			newpipeline.Target.RefType = ""
			newpipeline.Target.RefName = ""
		}

		pipeline := api.RunPipeline(repo, newpipeline)

		if pipeline.State.Result.Name == "" {
			fmt.Printf("%s", util.FormatPipelineStatus(pipeline.State.Name))
		} else {
			fmt.Printf("%s", util.FormatPipelineStatus(pipeline.State.Result.Name))
		}
		fmt.Printf(" \033[1;32m#%d\033[m ", pipeline.BuildNumber)
		if pipeline.Target.Source != "" {
			fmt.Printf("%s \033[1;34m[ %s → %s]\033[m\n", pipeline.Target.PullRequest.Title, pipeline.Target.Source, pipeline.Target.Destination)
		} else {
			fmt.Printf("\033[1;34m[ %s ]\033[m\n", pipeline.Target.RefName)
		}

		fmt.Printf("        \033[33m%s\033[m \033[37mTrigger: %s\033[m\n", pipeline.Author.DisplayName, pipeline.Trigger.Name)

		stepsChannel := api.GetPipelineSteps(repo, fmt.Sprintf("%d", pipeline.BuildNumber))
		for _, step := range <-stepsChannel {
			fmt.Printf("%s %s\n", step.Name, util.FormatPipelineStatus(step.State.Name))
		}
	},
}

func init() {
	RunCmd.Flags().StringP("select", "s", "", "select which pipeline definition to run")
	RunCmd.Flags().StringP("commit", "c", "", "run pipeline on branch for specific commit")
	RunCmd.Flags().StringP("branch", "b", "", "run pipeline for specific branch")
	RunCmd.Flags().StringP("pull-request", "p", "", "run pipeline for a specific pull-request")
	// TODO can we choose which step to log ?
	// LogsCmd.Flags().BoolP("tail", "t", false, "tail logs of a running pipeline step")
}
