package pipeline

import (
	"bb/api"
	"bb/util"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ViewCmd = &cobra.Command{
	Use:   "view [ID]",
	Short: "View details of a pipeline",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")
		showCommands, _ := cmd.Flags().GetBool("commands")
		branch, _ := cmd.Flags().GetString("target")

		var id int
		var err error
		if len(args) == 0 {
			if branch == "" {
				branch, err = util.GetCurrentBranch()
				cobra.CheckErr(err)
			}
			// retrieve id of pr for current branch
			pipeline := <-api.GetPipelineList(repo, 1, branch)
			if pipeline.BuildNumber == 0 {
				cobra.CheckErr(fmt.Sprintf("No pipelines found for target branch: '%s'", branch))
			}
			id = pipeline.BuildNumber
		} else {
			id, err = strconv.Atoi(args[0])
			cobra.CheckErr(err)
		}

		// make the steps request so that it's ready to print later on
		stepsChannel := api.GetPipelineSteps(repo, fmt.Sprintf("%d", id))
		pipeline := <-api.GetPipeline(repo, fmt.Sprintf("%d", id))

		if pipeline.State.Result.Name == "" {
			fmt.Printf("%s", util.FormatPipelineStatus(pipeline.State.Name))
		} else {
			fmt.Printf("%s", util.FormatPipelineStatus(pipeline.State.Result.Name))
		}
		fmt.Printf(" \033[1;32m#%d\033[m ", pipeline.BuildNumber)
		if pipeline.Target.Source != "" {
			fmt.Printf("%s \033[1;34m[ %s → %s] \033[37m%s\033[m\n", pipeline.Target.PullRequest.Title, pipeline.Target.Source, pipeline.Target.Destination, util.TimeAgo(pipeline.CreatedOn))
		} else {
			fmt.Printf("\033[1;34m[ %s ]\033[m\n", pipeline.Target.RefName)
		}

		fmt.Printf("        \033[33m%s\033[m \033[37mTrigger: %s\033[m\n", pipeline.Author.DisplayName, pipeline.Trigger.Name)

		fmt.Println()
		for _, step := range <-stepsChannel {
			if step.State.Result.Name != "" {
				fmt.Printf("%s %s \033[37m%s\033[m", step.Name, util.FormatPipelineStatus(step.State.Result.Name), util.TimeDuration(time.Duration(step.DurationInSeconds*1000000000)))
			} else if step.State.Stage.Name != "" {
				fmt.Printf("%s %s \033[37m%s\033[m", step.Name, util.FormatPipelineStatus(step.State.Stage.Name), util.TimeDuration(time.Duration(step.DurationInSeconds*1000000000)))
			} else {
				fmt.Printf("%s %s \033[37m%s\033[m", step.Name, util.FormatPipelineStatus(step.State.Name), util.TimeDuration(time.Duration(step.DurationInSeconds*1000000000)))
			}
			fmt.Println()
			if showCommands {
				for _, command := range step.ScriptCommands {
					fmt.Printf("\t%s\n", command.Name)
				}
			}
		}

		web, _ := cmd.Flags().GetBool("web")
		if web {
			util.OpenInBrowser(api.BBBrowsePipelines(repo, id))
			return
		}

	},
}

func init() {
	ViewCmd.Flags().String("target", "", "filter by target branch.")
	ViewCmd.RegisterFlagCompletionFunc("target", util.BranchCompletion)
	ViewCmd.Flags().Bool("web", false, "open in the browser")
	ViewCmd.Flags().BoolP("commands", "c", false, "show step commands")
}
