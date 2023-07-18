package issue

import (
	"bb/api"
	"bb/util"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ViewCmd = &cobra.Command{
	Use:   "view [KEY]",
	Short: "View issue",
	Run: func(cmd *cobra.Command, args []string) {
		var key string
		if len(args) == 0 {
			key = util.GetCurrentBranch()
		} else {
			key = args[0]
		}
		issue := <-api.GetIssue(viper.GetString("repo"), key)
		fmt.Printf("%v\n", issue)
	},
}

func init() {
}
