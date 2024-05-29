package issue

import (
	"bb/api"
	"bb/util"
	"fmt"
	"regexp"
	"strings"

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

		var data = api.UpdateIssueRequestBody{}
		// data.Fields.Summary = "TASKS - Data recovery for existing tasks"
		// TODO fix this in your project
		data.Fields.TimeTracking = &api.TimeTracking{
			OriginalEstimate: strings.Join(args[1:], " "),
		}
		// data.Update.TimeTracking = []api.UpdateType[api.TimeTracking]{
		// 	api.UpdateType[api.TimeTracking]{},
		// }
		// data.Update.TimeTracking[0].Set = &api.TimeTracking{
		// 	OriginalEstimate: strings.Join(args[1:], " "),
		// }

		api.UpdateIssue(key, data)
		issue := <-api.GetIssue(key)

		timeSpent := "-"
		if issue.Fields.TimeTracking.TimeSpent != " " {
			timeSpent = issue.Fields.TimeTracking.TimeSpent
		}

		fmt.Println()
		util.Printf("%s \033[1;32m%s\033[m %s %s %s\n", util.FormatIssueStatus(issue.Fields.Status.Name), issue.Key, util.FormatIssueType(issue.Fields.Type.Name), issue.Fields.Summary, util.FormatIssuePriority(issue.Fields.Priority.Id))
		util.Printf("    Assigned: \033[1;33m%s\033[m -> Reporter: \033[1;36m%s\033[m \033[37m(%d comments)\n", issue.Fields.Assignee.DisplayName, issue.Fields.Reporter.DisplayName, issue.Fields.Comment.Total)
		util.Printf("    Time spent: \033[1;34m%s\033[m [ %s/%s ]\n", timeSpent, issue.Fields.TimeTracking.RemainingEstimate, issue.Fields.TimeTracking.OriginalEstimate)
		fmt.Println()

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
	EstimateCmd.Flags().BoolP("transition", "t", false, "Also prompt to perform a transition")
}
