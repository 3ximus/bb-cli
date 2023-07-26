package issue

import (
	"bb/api"
	"bb/util"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

var LogCmd = &cobra.Command{
	Use:   "log [KEY] [TIME...]",
	Short: "Log time for an issue",
	Long: `Log time for an issue.
	Time format "2h 30m", "1d 5m" ...`,
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

		api.PostWorklog(key, seconds)
		fmt.Printf("Logged time for %s +\033[1;32m%d\033[m\n", key, seconds)
	},
}
