package issue

import (
	"bb/api"
	"bb/util"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

var PriorityCmd = &cobra.Command{
	Use:   "priority [KEY] [PRIORITY]",
	Short: "Set ticket priority",
	Long:  "Set ticket priority",
	Args:  cobra.MaximumNArgs(2),
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

		var data = api.UpdateIssueRequestBody{}

		if len(args) == 2 {
			data.Fields.Priority.Id = args[1]
		} else {
			priorities := []string{"1", "2", "3", "4", "5"}
			prioritieNames := []string{"Highest", "High", "Medium", "Low", "Lowest"}
			optIndex := util.SelectFZF(priorities, "New Priority > ", func(i int) string {
				return fmt.Sprintf("%s %s", util.FormatIssuePriority(priorities[i]), prioritieNames[i])
			})
			if len(optIndex) > 0 {
				data.Fields.Priority.Id = priorities[optIndex[0]]
			} else {
				return
			}
		}

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
	},
}

func init() {
}
