package pr

import (
	"bb/api"
	"bb/util"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests from a repository",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		author, _ := cmd.Flags().GetString("author")
		search, _ := cmd.Flags().GetString("search")
		pages, _ := cmd.Flags().GetInt("pages")
		state, _ := cmd.Flags().GetString("state")
		status, _ := cmd.Flags().GetBool("status")
		prChannel := api.GetPrList(viper.GetString("repo"), strings.ToUpper(state), author, search, pages, status)

		fmt.Printf("\n  Pull Requests for \033[1;36m%s\033[m\n\n", viper.GetString("repo"))
		count := 0
		for pr := range prChannel {
			// if we didn't provide filter don't show the pr status
			fmt.Printf("%s \033[1;32m#%d\033[m %s  \033[1;34m[ %s → %s]\033[m\n", util.FormatPrState(pr.State), pr.ID, pr.Title, pr.Source.Branch.Name, pr.Destination.Branch.Name)
			fmt.Printf("%s\033[33m%s\033[m  \033[37mComments: %d\033[m\n", strings.Repeat(" ", len(util.FormatPrState(pr.State))-4), pr.Author.Nickname, pr.CommentCount)
			count++
		}
		if count == 0 {
			fmt.Printf("\n  No pull requests for \033[1;36m%s\033[m\n\n", viper.GetString("repo"))
		} else {
			fmt.Println()
		}
	},
}

func init() {
	ListCmd.Flags().StringP("author", "a", "", "filter by author nick name (full nickname is needed due to an API limitation from bitbucket)")
	ListCmd.Flags().StringP("search", "S", "", "search pull request with query")
	ListCmd.Flags().StringP("state", "s", string(api.OPEN), `filter by state. Default: "open"
	possible options: "open", "merged", "declined" or "superseded"`)
	ListCmd.RegisterFlagCompletionFunc("state", stateCompletion)
	ListCmd.Flags().IntP("pages", "p", 1, "number of pages with results to retrieve")
	ListCmd.Flags().Bool("status", false, "include status of each pull request on the result. (the result will be slower)")

	// TODO
	ListCmd.Flags().String("destination", "", "filter by destination branch. \033[31mNot implemented\033[m")
	ListCmd.Flags().String("source", "", "filter by source branch. \033[31mNot implemented\033[m")
}

func stateCompletion(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"open", "merged", "declined", "superseded"}, cobra.ShellCompDirectiveDefault
}
