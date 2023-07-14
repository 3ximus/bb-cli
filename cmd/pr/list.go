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
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		author, _ := cmd.Flags().GetString("author")
		search, _ := cmd.Flags().GetString("search")
		pages, _ := cmd.Flags().GetInt("pages")
		states, _ := cmd.Flags().GetStringArray("state")
		allStates, _ := cmd.Flags().GetBool("all")
		if allStates {
			states = []string{string(api.OPEN), string(api.MERGED), string(api.DECLINED), string(api.SUPERSEDED)}
		}
		status, _ := cmd.Flags().GetBool("status")
		source, _ := cmd.Flags().GetString("source")
		destination, _ := cmd.Flags().GetString("destination")

		prChannel := api.GetPrList(viper.GetString("repo"), states, author, search, source, destination, pages, status)

		fmt.Printf("\n  Pull Requests for \033[1;36m%s\033[m\n\n", viper.GetString("repo"))
		count := 0
		for pr := range prChannel {
			// if we didn't provide filter don't show the pr status
			fmt.Printf("%-19s \033[1;32m#%d\033[m %s  \033[1;34m[ %s → %s]\033[m\n", util.FormatPrState(pr.State), pr.ID, pr.Title, pr.Source.Branch.Name, pr.Destination.Branch.Name)
			fmt.Printf("             \033[33m%s\033[m  \033[37mComments: %d\033[m", pr.Author.Nickname, pr.CommentCount)
			if status {
				fmt.Printf(" %s", util.FormatPipelineState(pr.Status.State))
			}
			fmt.Println()
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
	ListCmd.Flags().StringArrayP("state", "s", []string{string(api.OPEN)}, `filter by state. Default: "open". Multiple of these options can be given
	possible options: "open", "merged", "declined" or "superseded"`)
	ListCmd.RegisterFlagCompletionFunc("state", stateCompletion)
	ListCmd.Flags().IntP("pages", "p", 1, "number of pages with results to retrieve")
	ListCmd.Flags().Bool("status", false, "include status of each pull request on the result. (the result will be slower)")
	ListCmd.Flags().Bool("all", false, "return pull request with all possible states.")
	ListCmd.Flags().String("destination", "", "filter by destination branch.")
	ListCmd.Flags().String("source", "", "filter by source branch.")
}

func stateCompletion(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"open", "merged", "declined", "superseded"}, cobra.ShellCompDirectiveDefault
}
