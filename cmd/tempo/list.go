package tempo

import (
	"bb/api"
	"bb/util"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:     "list [ PROJECT-KEY ]",
	Short:   "List worklogs",
	Aliases: []string{"ls"},
	Long: `List worklogs from Tempo with some filters
	By default it filters worklogs assigned to the current user for the current day.
	Given an argument it will filter worklogs from that project. Otherwise it will try to derive the project name from the branch name.
	`,
	Args:    cobra.MaximumNArgs(1),
	Example: "list  ",
	Run: func(cmd *cobra.Command, args []string) {
		user := api.GetMyself()

		// TODO DEFAULT - allow flags to control this
		// list today's worklogs
		now := time.Now().UTC()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999000000, time.UTC)

		for _, w := range api.ListWorklogs(user, start, end) {
			startTime, err := time.Parse(time.RFC3339, w.StartDateTimeUtc)
			cobra.CheckErr(err)
			issue := <-api.GetIssue(strconv.Itoa(w.Issue.ID))

			util.Printf("\033[1;34m%s\033[m +\033[1;32m%s\033[m - \033[1;33m%s\033[m %s\n", startTime.Local().Format("15:00"), util.TimeDuration(time.Duration(w.TimeSpentSeconds*1e9)), issue.Key, issue.Fields.Summary)
		}
	},
}

func init() {
}
