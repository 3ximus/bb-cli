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
	Long:  "List issues from Jira with preset filtering. By default it filters tickets assigned to the current user",
	Run: func(cmd *cobra.Command, args []string) {
		nResults, _ := cmd.Flags().GetInt("results")
		reporter, _ := cmd.Flags().GetBool("reporter")
		all, _ := cmd.Flags().GetBool("all")
		statuses, _ := cmd.Flags().GetStringArray("status")
		project, _ := cmd.Flags().GetString("project")

		fmt.Println()
		for issue := range api.GetIssueList(viper.GetString("repo"), nResults, all, reporter, project, statuses) {
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
	ListCmd.Flags().BoolP("all", "a", false, "filter all issues. (Not assigned or reporting to current user)")
	ListCmd.Flags().BoolP("reporter", "r", false, "filter issues reporting to current user")
	ListCmd.Flags().IntP("results", "n", 10, "max number of results retrieve")
	ListCmd.Flags().StringArrayP("status", "s", []string{}, `filter status type.
	possible options: "todo", "inprogress", "testing", "done", "blocked"`)
	ListCmd.RegisterFlagCompletionFunc("status", statusCompletion)
}

func statusCompletion(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"todo", "inprogress", "testing", "done", "blocked"}, cobra.ShellCompDirectiveDefault
}
