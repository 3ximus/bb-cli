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

		api.PostWorklog(key, seconds)
		fmt.Printf("Logged time for %s +\033[1;32m%d\033[m\n", key, seconds)

		if transition {
			// select new state
			var newState = ""
			transitions := <-api.GetTransitions(key)
			var newStateName = ""
			optIndex := util.SelectFZF(transitions, "Transition To > ", func(i int) string {
				return fmt.Sprintf("%s", transitions[i].To.Name)
			})
			if len(optIndex) > 0 {
				newState = transitions[optIndex[0]].Id
				newStateName = transitions[optIndex[0]].To.Name
			}
			if key == "" || newState == "" {
				return
			}

			api.PostTransitions(key, newState)
			fmt.Printf("Issue status changed for %s -> \033[1;32m%s\033[m\n", key, newStateName)
		}

	},
}

func init() {
	LogCmd.Flags().BoolP("transition", "t", false, "Also prompt to perform a transition")
}
