package issue

import (
	"bb/api"
	"bb/util"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

var EstimateCmd = &cobra.Command{
	Use:   "estimate [KEY] [TIME...]",
	Short: "Estimate time for an issue",
	Long: `Estimate time for an issue.
	Time format "2h 30m", "1d 5m" ...`,
	ValidArgsFunction: func(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return ListBranchesMatchingJiraTickets(), cobra.ShellCompDirectiveDefault
		} else {
			return []string{}, cobra.ShellCompDirectiveDefault
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		transition, _ := cmd.Flags().GetBool("transition")

		var key string
		if len(args) == 0 {
			branch, err := util.GetCurrentBranch()
			cobra.CheckErr(err)
			re := regexp.MustCompile(api.JiraIssueKeyRegex)
			key = re.FindString(branch)
			// TODO maybe use an option to get the key from a PR ?
		} else {
			key = args[0]
		}

		seconds, err := util.ConvertToSeconds(args[1:])
		cobra.CheckErr(err)

		fmt.Println("Not implemented.", key, seconds)

		if transition {
		}
	},
}

func init() {
	EstimateCmd.Flags().BoolP("transition", "t", false, "Also prompt to perform a transition")
}
