package issue

import (
	"bb/api"
	"bb/util"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

var LogCmd = &cobra.Command{
	Use:     "log [KEY] [TIME]",
	Short:   "Log time for an issue",
	Args:    cobra.MaximumNArgs(2),
	ValidArgsFunction: func(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return ListBranchesMatchingJiraTickets(), cobra.ShellCompDirectiveDefault
		} else {
			return []string{}, cobra.ShellCompDirectiveDefault
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var key string
		if len(args) == 0 {
			branch := util.GetCurrentBranch()
			re := regexp.MustCompile(api.JiraIssueKeyRegex)
			key = re.FindString(branch)
			// TODO maybe use an option to get the key from a PR ?
		} else {
			key = args[0]
		}

		fmt.Println("Not implemented.", key)
	},
}
