package issue

import (
	"bb/api"
	"bb/util"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:   "list [ PROJECT-KEY | all ]",
	Short: "List issues",
	Long: `List issues from Jira with preset filtering.
	By default it filters tickets assigned to the current user and it tries to gess the current project from the current branch name.
	Given an argument it will filter tickets from that project. Otherwise it will try to derive the project name from the branch name.
	If all is given then project filtering is not applied
	`,
	Args: cobra.MaximumNArgs(1),
	Example: "list DP --search ",
	Run: func(cmd *cobra.Command, args []string) {
		nResults, _ := cmd.Flags().GetInt("number-results")
		reporter, _ := cmd.Flags().GetBool("reporter")
		all, _ := cmd.Flags().GetBool("all")
		search, _ := cmd.Flags().GetString("search")
		statuses, _ := cmd.Flags().GetStringArray("status")
		priority, _ := cmd.Flags().GetBool("priority")
		showUsers, _ := cmd.Flags().GetBool("users")
		showTime, _ := cmd.Flags().GetBool("time")

		project := ""
		if len(args) == 0 {
			branch, err := util.GetCurrentBranch()
			if err == nil {
				re := regexp.MustCompile(api.JiraIssueKeyRegex)
				key := re.FindString(branch)
				if key != "" {
					project = strings.Split(key, "-")[0]
				}
			}
		} else {
			project = args[0]
			if project == "all" {
				project = ""
			}
		}

		// convert status based on current settings
		var statusConversion = []string{}
		for _, s := range statuses {
			statusMap, err := util.GetConfig("jira_status", s)
			if err == nil {
				statusConversion = append(statusConversion, statusMap.Values...)
			} else {
				statusConversion = append(statusConversion, s)
			}
		}

		for issue := range api.GetIssueList(nResults, all, reporter, project, statusConversion, search, priority) {
			timeSpent := "-"
			if issue.Fields.TimeTracking.TimeSpent != " " {
				timeSpent = issue.Fields.TimeTracking.TimeSpent
			}
			util.Printf("%s \033[1;32m%s\033[m %s %s %s\n", util.FormatIssueStatus(issue.Fields.Status.Name), issue.Key, util.FormatIssueType(issue.Fields.Type.Name), issue.Fields.Summary, util.FormatIssuePriority(issue.Fields.Priority.Id, issue.Fields.Priority.Name))
			if showUsers {
				if all {
					util.Printf("    \033[37mAssigned: \033[1m%s\033[0;37m -> Reporter: \033[1;36m%s\033[m \033[37m(%d comments)\033[m\n", issue.Fields.Assignee.DisplayName, issue.Fields.Reporter.DisplayName, issue.Fields.Comment.Total)
				} else if reporter {
					util.Printf("    \033[37mAssigned: \033[1m%s \033[37m(%d comments)\033[m\n", issue.Fields.Assignee.DisplayName, issue.Fields.Comment.Total)
				} else {
					util.Printf("    \033[37mReporter: \033[1m%s \033[37m(%d comments)\033[m\n", issue.Fields.Reporter.DisplayName, issue.Fields.Comment.Total)
				}
			}
			if showTime {
				util.Printf("    \033[37mTime spent: \033[1;34m%s\033[m [ %s/%s ]\n", timeSpent, issue.Fields.TimeTracking.RemainingEstimate, issue.Fields.TimeTracking.OriginalEstimate)
			}
		}
	},
}

func init() {
	// filter
	ListCmd.Flags().BoolP("all", "a", false, "filter all issues. (Not assigned or reporting to current user)")
	ListCmd.Flags().BoolP("reporter", "r", false, "filter issues reporting to current user")
	ListCmd.Flags().StringArrayP("status", "s", []string{}, `filter status type.
	this option will provide completion for the mappings defined in "jira_status" of your config file`)
	ListCmd.RegisterFlagCompletionFunc("status", statusCompletion)
	ListCmd.Flags().String("search", "", "filter issues by keyword")

	// TODO add way to sort by recent or the ones the user has participated on

	// display
	ListCmd.Flags().BoolP("users", "u", false, "show users")
	ListCmd.Flags().BoolP("time", "t", false, "show time information")
	ListCmd.Flags().IntP("number-results", "n", 10, "max number of results retrieve")
	// sort
	ListCmd.Flags().BoolP("priority", "P", false, "sort by priority")
}

func statusCompletion(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	statusMap := viper.GetStringMapStringSlice("jira_status")
	status := make([]string, 0, len(statusMap))
	for k := range statusMap {
		status = append(status, k)
	}
	return status, cobra.ShellCompDirectiveDefault
}
