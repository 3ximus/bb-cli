package issue

import (
	"bb/api"
	"bb/util"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

var TransitionCmd = &cobra.Command{
	Use:     "transition [KEY]...",
	Short:   "Transition issues to another state",
	Long:    "Transition issues to another state. You'll be prompted to choose one of the available states",
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"t"},
	ValidArgsFunction: func(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return ListBranchesMatchingJiraTickets(), cobra.ShellCompDirectiveDefault
		} else {
			return []string{}, cobra.ShellCompDirectiveDefault
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var keys []string
		if len(args) == 0 {
			branch, err := util.GetCurrentBranch()
			cobra.CheckErr(err)
			re := regexp.MustCompile(api.JiraIssueKeyRegex)
			keys = []string{re.FindString(branch)}
			// TODO maybe use an option to get the key from a PR
		} else {
			keys = args
		}

		for _, key := range keys {
			// select new state
			var newState = ""
			transitions := <-api.GetTransitions(key)
			var newStateName = ""
			optIndex := util.SelectFZF(transitions, fmt.Sprintf("Transition %s To > ", key), func(i int) string {
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
}

func ListBranchesMatchingJiraTickets() []string {
	var branches = []string{}
	re := regexp.MustCompile(api.JiraIssueKeyRegex)
	for _, branch := range util.ListBranches() {
		if re.MatchString(branch) {
			branches = append(branches, re.FindString(branch))
		}
	}
	return branches
}
