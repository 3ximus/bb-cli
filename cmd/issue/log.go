package issue

import (
	"bb/api"
	"bb/util"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var LogCmd = &cobra.Command{
	Use:   "log [KEY] [TIME...]",
	Short: "Log time for an issue",
	Long: `Log time for an issue.
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

		var seconds int
		var err error
		if len(args) == 2 {
			seconds, err = util.ConvertToSeconds(args[1:])
		} else {
			fmt.Print("? \033[1;35mLog time \033[m")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			seconds, err = util.ConvertToSeconds(strings.Split(scanner.Text(), " "))
		}
		cobra.CheckErr(err)

		user := api.GetMyself()
		issueChan := api.GetIssue(key)

		// list today's worklogs
		now := time.Now().UTC()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999000000, time.UTC)

		timeStartWorklog := time.Date(now.Year(), now.Month(), now.Day(), viper.GetInt("day_start_hour"), 0, 0, 0, time.Local)
		for _, w := range api.ListWorklogs(user, start, end) {
			startTime, err := time.Parse(time.RFC3339, w.StartDateTimeUtc)
			cobra.CheckErr(err)
			lastWorklog := startTime.Local().Add(time.Duration(w.TimeSpentSeconds * 1e9))
			if lastWorklog.After(timeStartWorklog) {
				timeStartWorklog = lastWorklog
			}
		}
		issueId, err := strconv.Atoi((<-issueChan).ID)
		cobra.CheckErr(err)
		newWorklog := api.PostWorklog(user, issueId, seconds, timeStartWorklog)
		newStartTime, err := time.Parse(time.RFC3339, newWorklog.StartDateTimeUtc)
		cobra.CheckErr(err)
		fmt.Printf("Logged time for %s  |  \033[1;34m%s\033[m +\033[1;32m%s\033[m\n", key, newStartTime.Local().Format("15:04"), util.TimeDuration(time.Duration(newWorklog.TimeSpentSeconds*1e9)))

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
	LogCmd.Flags().BoolP("transition", "t", false, "Also prompt to perform a transition")
}
