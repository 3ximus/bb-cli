package pipeline

import (
	"bb/api"
	"bb/util"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var LogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show logs of a pipeline",
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

		var stepUUID = ""
		steps := <-api.GetPipelineSteps(repo, fmt.Sprintf("%d", id))
		selectedStep, _ := cmd.Flags().GetString("step")
		if selectedStep == "" {
			optIndex := util.SelectFZF(steps, fmt.Sprintf("Step to Log > "), func(i int) string {
				return fmt.Sprintf("%s", steps[i].Name)
			})
			if len(optIndex) > 0 {
				stepUUID = steps[optIndex[0]].UUID
			}
		} else {
			for _, step := range steps {
				if step.Name == selectedStep || strings.ToLower(step.Name) == selectedStep {
					stepUUID = step.UUID
				}
			}
			if stepUUID == "" {
				cobra.CheckErr("Step not found")
			}
		}
		fmt.Print(<-api.GetPipelineStepLogs(repo, fmt.Sprintf("%d", id), stepUUID))
	},
}

func init() {
	LogsCmd.Flags().StringP("step", "s", "", "select step. Without this option the step is prompet interactively")
}
