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
		nResults, _ := cmd.Flags().GetInt("results")
		assignee, _ := cmd.Flags().GetBool("assignee")
		reporter, _ := cmd.Flags().GetBool("reporter")
		project, _ := cmd.Flags().GetString("project")

		fmt.Println()
		for issue := range api.GetIssueList(viper.GetString("repo"), nResults, assignee, reporter, project) {
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
	ListCmd.Flags().StringP("project", "p", "", "filter issues by project key")
	ListCmd.Flags().BoolP("assignee", "a", true, "filter issues assigned to me. Default operation")
	ListCmd.Flags().BoolP("reporter", "r", false, "filter issues reporting to me")
	ListCmd.Flags().IntP("results", "n", 10, "max number of results retrieve")
}
