package pipeline

import (
	"bb/api"
	"bb/util"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var LogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show logs of a pipeline step",
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

		var selected = api.PipelineStep{}
		steps := <-api.GetPipelineSteps(repo, fmt.Sprintf("%d", id))
		selectedStep, _ := cmd.Flags().GetString("step")
		if selectedStep == "" {
			optIndex := util.SelectFZF(steps, fmt.Sprintf("Step to Log > "), func(i int) string {
				return fmt.Sprintf("%s", steps[i].Name)
			})
			if len(optIndex) > 0 {
				selected = steps[optIndex[0]]
			}
		} else {
			for _, step := range steps {
				if step.Name == selectedStep || strings.ToLower(step.Name) == selectedStep {
					selected = step
				}
			}
		}

		if selected.UUID == "" {
			cobra.CheckErr("Step not found")
		}

		tail, _ := cmd.Flags().GetBool("tail")
		if !tail {
			fmt.Print(<-api.GetPipelineStepLogs(repo, fmt.Sprintf("%d", id), selected.UUID, 0))
		} else {
			firstDone := false
			totalLength := 0
			for !firstDone || selected.State.Name != "COMPLETED" {
				if firstDone {
					time.Sleep(2 * time.Second)
				}
				logsChannel := api.GetPipelineStepLogs(repo, fmt.Sprintf("%d", id), selected.UUID, totalLength)
				stepChannel := api.GetPipelineStep(repo, fmt.Sprintf("%d", id), selected.UUID)
				response := <-logsChannel
				fmt.Print(response)
				totalLength += len(response)
				selected = <-stepChannel
				firstDone = true
			}
		}
	},
}

func init() {
	LogsCmd.Flags().StringP("step", "s", "", "select step. Without this option the step is prompet interactively")
	LogsCmd.Flags().BoolP("tail", "t", false, "tail logs of a running pipeline step")
}
