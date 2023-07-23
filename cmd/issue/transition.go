package issue

import (
	"bb/api"
	"bb/util"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var TransitionCmd = &cobra.Command{
	Use:   "transition [KEY]",
	Short: "Transition issue to another state",
	Run: func(cmd *cobra.Command, args []string) {
		var key string
		if len(args) == 0 {
			branch := util.GetCurrentBranch()
			re := regexp.MustCompile(api.JiraIssueKeyRegex)
			key = re.FindString(branch)
			// TODO maybe use an option to get the key from a PR
		} else {
			key = args[0]
		}
		transitions := <-api.GetTransitions(viper.GetString("repo"), key)

		fmt.Println(transitions)

	},
}

func init() {
}
