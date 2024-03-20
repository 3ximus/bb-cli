package issue

import (
	"bb/api"
	"bb/util"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ViewCmd = &cobra.Command{
	Use:   "view [KEY]",
	Short: "View issue",
	Run: func(cmd *cobra.Command, args []string) {
		var key string
		if len(args) == 0 {
			branch, err := util.GetCurrentBranch()
			cobra.CheckErr(err)
			re := regexp.MustCompile(api.JiraIssueKeyRegex)
			key = re.FindString(branch)
			// TODO maybe use an option to get the key from a PR
		} else {
			key = args[0]
		}
		issue := <-api.GetIssue(key)

		timeSpent := "-"
		if issue.Fields.TimeTracking.TimeSpent != " " {
			timeSpent = issue.Fields.TimeTracking.TimeSpent
		}

		fmt.Println()
		util.Printf("%s \033[1;32m%s\033[m %s %s %s\n", util.FormatIssueStatus(issue.Fields.Status.Name), issue.Key, util.FormatIssueType(issue.Fields.Type.Name), issue.Fields.Summary, util.FormatIssuePriority(issue.Fields.Priority.Id, issue.Fields.Priority.Name))
		util.Printf("    Assigned: \033[1;33m%s\033[m -> Reporter: \033[1;36m%s\033[m \033[37m(%d comments)\n", issue.Fields.Assignee.DisplayName, issue.Fields.Reporter.DisplayName, issue.Fields.Comment.Total)
		util.Printf("    Time spent: \033[1;34m%s\033[m [ %s/%s ]\n", timeSpent, issue.Fields.TimeTracking.RemainingEstimate, issue.Fields.TimeTracking.OriginalEstimate)
		if issue.Fields.Parent.Fields.Summary != "" {
			util.Printf("    \033[37mParent:\033[m %s %s (\033[37m%s\033[m)\n", util.FormatIssueType(issue.Fields.Parent.Fields.Type.Name), issue.Fields.Parent.Fields.Summary, issue.Fields.Parent.Key)
		} else {
			util.Printf("    \033[37mParent: ---\n")
		}
		fmt.Println()

		web, _ := cmd.Flags().GetBool("web")
		if web {
			util.OpenInBrowser(api.JiraBrowse(viper.GetString("jira_domain"), key))
			return
		}
	},
}

func init() {
	ViewCmd.Flags().Bool("web", false, "Open in the browser.")
}
