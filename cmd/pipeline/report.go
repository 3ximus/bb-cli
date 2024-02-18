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

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Show test reports of a pipeline step",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")

		var id int
		var err error
		if len(args) == 0 {
			branch, err := util.GetCurrentBranch()
			cobra.CheckErr(err)
			// retrieve id of pipeline for current branch
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

		fullReportChannel := api.GetPipelineReportCases(repo, fmt.Sprintf("%d", id), selected.UUID)
		report := <-api.GetPipelineReport(repo, fmt.Sprintf("%d", id), selected.UUID)
		fmt.Println("Test report:")
		fmt.Printf("\033[1;32mPassed:  %3d\033[m\n", report.Success)
		fmt.Printf("\033[1;31mFailed:  %3d\033[m\n", report.Failed)
		fmt.Printf("\033[1;33mSkipped: %3d\033[m\n", report.Skipped)
		fmt.Printf("Error:   %3d\n", report.Error)
		fmt.Printf("Total:   %3d\n", report.Total)

		showFull, _ := cmd.Flags().GetBool("full")
		if showFull {
			for reportCase := range fullReportChannel {
				fmt.Printf("%s \033[34m%s\033[m %s \033[37m%s\033[m\n", util.FormatPipelineStatus(reportCase.Status), reportCase.PackageName, reportCase.Name, strings.Replace(reportCase.Duration, "PT", "", 1))
			}
		}
	},
}

func init() {
	ReportCmd.Flags().StringP("step", "s", "", "select step. Without this option the step is prompet interactively")
	ReportCmd.Flags().BoolP("full", "f", false, "show the full report")
}
