package issue

import (
	"bb/api"
	"bb/util"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO add a default status listing on the config file

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	Long: `List issues from Jira with preset filtering.
	By default it filters tickets assigned to the current user and it tries to gess the current project from the current branch name.`,
	Run: func(cmd *cobra.Command, args []string) {
		nResults, _ := cmd.Flags().GetInt("number-results")
		reporter, _ := cmd.Flags().GetBool("reporter")
		all, _ := cmd.Flags().GetBool("all")
		statuses, _ := cmd.Flags().GetStringArray("status")
		project, _ := cmd.Flags().GetString("project")
		priority, _ := cmd.Flags().GetBool("priority")
		showUsers, _ := cmd.Flags().GetBool("show-users")

		if !cmd.Flags().Changed("project") {
			branch, err := util.GetCurrentBranch()
			if err == nil {
				re := regexp.MustCompile(api.JiraIssueKeyRegex)
				key := re.FindString(branch)
				if key != "" {
					project = strings.Split(key, "-")[0]
				}
			}
		}
		if project == "all" {
			project = ""
		}

		// convert status based on current settings
		var statusConversion = []string{}
		statusMap := viper.GetStringMapStringSlice("jira_status")
		for _, s := range statuses {
			if val, exists := statusMap[s]; exists {
				for _, k := range val {
					statusConversion = append(statusConversion, k)
				}
			} else {
				statusConversion = append(statusConversion, s)
			}
		}

		fmt.Println()
		for issue := range api.GetIssueList(nResults, all, reporter, project, statusConversion, priority) {
			timeSpent := "-"
			if issue.Fields.TimeTracking.TimeSpent != " " {
				timeSpent = issue.Fields.TimeTracking.TimeSpent
			}
			fmt.Printf("%s \033[1;32m%s\033[m %s %s %s\n", util.FormatIssueStatus(issue.Fields.Status.Name), issue.Key, util.FormatIssueType(issue.Fields.Type.Name), issue.Fields.Summary, util.FormatIssuePriority(issue.Fields.Priority.Id, issue.Fields.Priority.Name))
			// TODO format spacing better
			if showUsers {
				if all {
					fmt.Printf("    \033[37mAssigned: \033[1m%s\033[0;37m -> Reporter: \033[1;36m%s\033[m \033[37m(%d comments)\033[m\n", issue.Fields.Assignee.DisplayName, issue.Fields.Reporter.DisplayName, issue.Fields.Comment.Total)
				} else if reporter {
					fmt.Printf("    \033[37mAssigned: \033[1m%s \033[37m(%d comments)\033[m\n", issue.Fields.Assignee.DisplayName, issue.Fields.Comment.Total)
				} else {
					fmt.Printf("    \033[37mReporter: \033[1m%s \033[37m(%d comments)\033[m\n", issue.Fields.Reporter.DisplayName, issue.Fields.Comment.Total)
				}
			}
			fmt.Printf("    \033[37mTime spent: \033[1;34m%s\033[m [ %s/%s ]\n", timeSpent, issue.Fields.TimeTracking.RemainingEstimate, issue.Fields.TimeTracking.OriginalEstimate)
		}
		fmt.Println()
	},
}

func init() {
	// filters
	ListCmd.Flags().StringP("project", "p", "", `filter issues by project key.
	If "all" is given it shows all projects (when a project is detected on current branch and you still want to show all projects)`)
	ListCmd.Flags().BoolP("all", "a", false, "filter all issues. (Not assigned or reporting to current user)")
	ListCmd.Flags().BoolP("reporter", "r", false, "filter issues reporting to current user")
	ListCmd.Flags().StringArrayP("status", "s", []string{}, `filter status type.
	possible options: "open", "todo", "inprogress", "testing", "done", "blocked"`)
	ListCmd.RegisterFlagCompletionFunc("status", statusCompletion)

	// TODO add way to sort by recent or the ones the user has participated on

	// display
	ListCmd.Flags().BoolP("show-users", "u", false, "show users")
	ListCmd.Flags().BoolP("one-line", "o", false, "show one issue per line. \033[31mNot implemented\033[m")
	ListCmd.Flags().IntP("number-results", "n", 10, "max number of results retrieve")
	// sort
	ListCmd.Flags().BoolP("priority", "P", false, "sort by priority")
}

func statusCompletion(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"open", "todo", "inprogress", "testing", "done", "blocked"}, cobra.ShellCompDirectiveDefault
}
