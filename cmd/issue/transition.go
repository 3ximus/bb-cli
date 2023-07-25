package issue

import (
	"bb/api"
	"bb/util"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

var TransitionCmd = &cobra.Command{
	Use:     "transition [KEY] [NEW STATE]",
	Short:   "Transition issue to another state",
	Long:    "Transition issue to another state. If no state is given you'll be prompted to choose on of the available states",
	Args:    cobra.MaximumNArgs(2),
	Aliases: []string{"t"},
	ValidArgsFunction: func(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return ListBranchesMatchingJiraTickets(), cobra.ShellCompDirectiveDefault
		} else if len(args) == 1 {
			transitions := <-api.GetTransitions(args[0])
			// TODO fix this completion because the options are split into multiple strings by the space
			var opt = []string{}
			for _, t := range transitions {
				opt = append(opt, "\""+t.To.Name+"\"")
			}
			return opt, cobra.ShellCompDirectiveDefault
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
			// TODO maybe use an option to get the key from a PR
		} else {
			key = args[0]
		}

		// select new state
		var newState = ""
		transitions := <-api.GetTransitions(key)
		var newStateName = ""
		if len(args) == 1 {
			optIndex := util.UseExternalFZF(transitions, "Transition To > ", func(i int) string {
				return fmt.Sprintf("%s", transitions[i].To.Name)
			})
			if len(optIndex) > 0 {
				newState = transitions[optIndex[0]].Id
				newStateName = transitions[optIndex[0]].To.Name
			}
		} else {
			for _, t := range transitions {
				if t.To.Name == args[1] {
					newState = t.Id
					newStateName = t.To.Name
					break
				}
			}
		}
		if key == "" || newState == "" {
			return
		}

		api.PostTransitions(key, newState)
		fmt.Printf("Issue status changed for %s -> \033[1;32m%s\033[m\n", key, newStateName)
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
