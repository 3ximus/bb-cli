package issue

import (
	"bb/api"
	"bb/util"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:     "list [ PROJECT-KEY ]",
	Short:   "List issues",
	Aliases: []string{"ls"},
	Long: `List issues from Jira with preset filtering.
	By default it filters tickets assigned to the current user and it tries to gess the current project from the current branch name.
	Given an argument it will filter tickets from that project. Otherwise it will try to derive the project name from the branch name.
	If all is given then project filtering is not applied
	`,
	Args:    cobra.MaximumNArgs(1),
	Example: "list DP --search ",
	Run: func(cmd *cobra.Command, args []string) {
		nResults, _ := cmd.Flags().GetInt("number-results")
		reporter, _ := cmd.Flags().GetBool("reporter")
		all, _ := cmd.Flags().GetBool("all")
		search, _ := cmd.Flags().GetString("search")
		statuses, _ := cmd.Flags().GetStringArray("status")
		iTypes, _ := cmd.Flags().GetStringArray("type")
		priority, _ := cmd.Flags().GetBool("priority")
		lastWorked, _ := cmd.Flags().GetBool("last")
		showUsers, _ := cmd.Flags().GetBool("users")
		showTime, _ := cmd.Flags().GetBool("time")
		showParents, _ := cmd.Flags().GetBool("parent")

		useFZF, _ := cmd.Flags().GetBool("fzf")
		useFZFInternal, _ := cmd.Flags().GetBool("fzf-internal")
		if useFZF {
			util.ReplaceListWithFzf("--read0 --prompt 'View > '" +
				" --header='\033[1;33mctrl-w\033[m: web view | \033[1;33mctrl-t\033[m: transition issue | \033[1;33mctrl-p\033[m: set priority | \033[1;33mctrl-l\033[m: log time'" +
				" --preview '" + os.Args[0] + " issue view {2} --color' --preview-window=hidden" + // start with hidden preview
				" --bind 'enter:become(" + os.Args[0] + " issue view {2})'" +
				" --bind 'ctrl-w:execute(" + os.Args[0] + " issue view --web {2})'" +
				" --bind 'ctrl-t:execute(" + os.Args[0] + " issue transition {2})'" +
				" --bind 'ctrl-l:execute(" + os.Args[0] + " issue log {2})'" +
				" --bind 'ctrl-p:execute(" + os.Args[0] + " issue priority {2})'")
			return
		}

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

		// convert type based on current settings
		var typeConversion = []string{}
		for _, s := range iTypes {
			typesMap, err := util.GetConfig("jira_type", s)
			if err == nil {
				typeConversion = append(typeConversion, typesMap.Values...)
			} else {
				typeConversion = append(typeConversion, s)
			}
		}

		for issue := range api.GetIssueList(nResults, all, reporter, project, statusConversion, typeConversion, search, priority, lastWorked) {
			timeSpent := "-"
			if issue.Fields.TimeTracking.TimeSpent != " " {
				timeSpent = issue.Fields.TimeTracking.TimeSpent
			}
			util.Printf("%s \033[1;32m%s\033[m %s %s %s", util.FormatIssueStatus(issue.Fields.Status.Name), issue.Key, util.FormatIssueType(issue.Fields.Type.Name), issue.Fields.Summary, util.FormatIssuePriority(issue.Fields.Priority.Id))
			if showParents {
				if issue.Fields.Parent.Fields.Summary != "" {
					util.Printf("\n    \033[37mParent:\033[m %s %s (\033[37m%s\033[m)", util.FormatIssueType(issue.Fields.Parent.Fields.Type.Name), issue.Fields.Parent.Fields.Summary, issue.Fields.Parent.Key)
				} else {
					util.Printf("\n    \033[37mParent: ---")
				}
			}
			if showUsers {
				if all {
					util.Printf("\n    \033[37mAssigned: \033[1m%s\033[0;37m -> Reporter: \033[1;36m%s\033[m \033[m\033[37m(%d comments)\033[m", issue.Fields.Assignee.DisplayName, issue.Fields.Reporter.DisplayName, issue.Fields.Comment.Total)
				} else if reporter {
					util.Printf("\n    \033[37mAssigned: \033[1m%s \033[m\033[37m(%d comments)\033[m", issue.Fields.Assignee.DisplayName, issue.Fields.Comment.Total)
				} else {
					util.Printf("\n    \033[37mReporter: \033[1m%s \033[m\033[37m(%d comments)\033[m", issue.Fields.Reporter.DisplayName, issue.Fields.Comment.Total)
				}
			}
			if showTime {
				util.Printf("\n    \033[37mTime spent: \033[1;34m%s\033[m [ %s/%s ]", timeSpent, issue.Fields.TimeTracking.RemainingEstimate, issue.Fields.TimeTracking.OriginalEstimate)
			}

			endChar := "\n"
			if useFZFInternal {
				endChar = "\x00"
			}
			util.Printf(endChar)
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
	ListCmd.Flags().StringArray("type", []string{}, `filter issue type.
	this option will provide completion for the mappings defined in "jira_type" of your config file`)
	ListCmd.RegisterFlagCompletionFunc("type", typeCompletion)
	ListCmd.Flags().String("search", "", "filter issues by keyword")

	// display
	ListCmd.Flags().BoolP("users", "u", false, "show users")
	ListCmd.Flags().BoolP("time", "t", false, "show time information")
	ListCmd.Flags().BoolP("parent", "p", false, "show parent tickets")
	ListCmd.Flags().IntP("number-results", "n", 99, "max number of results retrieve")
	ListCmd.Flags().BoolP("last", "l", false, "display tickets where status was changed by current user. sorted by last updated")
	// sort
	ListCmd.Flags().BoolP("priority", "P", false, "sort by priority")

	if util.CommandExists("fzf") {
		ListCmd.Flags().Bool("fzf", false, "use fzf interface on results")
		ListCmd.Flags().Bool("fzf-internal", false, "use fzf interface on results")
		ListCmd.Flags().MarkHidden("fzf-internal")
	}
}

func statusCompletion(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	statusMap := viper.GetStringMapStringSlice("jira_status")
	status := make([]string, 0, len(statusMap))
	for k := range statusMap {
		status = append(status, k)
	}
	return status, cobra.ShellCompDirectiveDefault
}

func typeCompletion(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	issueTypeMap := viper.GetStringMapStringSlice("jira_type")
	issueType := make([]string, 0, len(issueTypeMap))
	for k := range issueTypeMap {
		issueType = append(issueType, k)
	}
	return issueType, cobra.ShellCompDirectiveDefault
}
