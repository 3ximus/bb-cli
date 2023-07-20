package issue

import (
	"bb/api"
	"bb/util"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println()
		for issue := range api.GetIssueList(viper.GetString("repo"), 1) {
			timeSpent := "-"
			if issue.Fields.TimeTracking.TimeSpent == " " {
				timeSpent = issue.Fields.TimeTracking.TimeSpent
			}
			fmt.Printf("%s \033[1;32m%s\033[m %s %s\n", util.FormatIssueStatus(issue.Fields.Status.Name), issue.Key, issue.Fields.Summary, util.FormatIssuePriority(issue.Fields.Priority.Id, issue.Fields.Priority.Name))
			fmt.Printf("    Assigned: \033[1;33m%s\033[m -> Reporter: \033[1;36m%s\033[m \033[37m(%d comments)\n", issue.Fields.Assignee.DisplayName, issue.Fields.Reporter.DisplayName, issue.Fields.Comment.Total)
			fmt.Printf("    Time spent: \033[1;34m%s\033[m [ %s/%s ]\n", timeSpent, issue.Fields.TimeTracking.OriginalEstimate, issue.Fields.TimeTracking.RemainingEstimate)
		}
		fmt.Println()
	},
}

func init() {
}
