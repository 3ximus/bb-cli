package pr

import (
	"bb/api"
	"bb/util"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests from a repository",
	Run: func(cmd *cobra.Command, args []string) {
		prs := api.GetPr(viper.GetString("repository"), []string{"OPEN"})
		prsdata := make([][]string, len(prs))
		for i, pr := range prs {
			prsdata[i] = []string{fmt.Sprint(pr.ID), pr.Title, pr.State, fmt.Sprint(pr.CommentCount)}
		}
		util.Table([]string{"ID", "Title", "State", "CommentCount"}, prsdata)
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
