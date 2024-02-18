package pr

import (
	"bb/api"
	"bb/util"
	"fmt"
	"strings"

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
		source, _ := cmd.Flags().GetString("source")
		destination, _ := cmd.Flags().GetString("destination")
		status, _ := cmd.Flags().GetBool("status")
		participants, _ := cmd.Flags().GetBool("participants")

		count := 0
		for pr := range api.GetPrList(viper.GetString("repo"), states, author, search, source, destination, pages, status, participants) {
			fmt.Printf("%s \033[1;32m#%d\033[m %s \033[1;34m[ %s \033[m→\033[1;34m %s ]\033[m \033[33m%s\033[m", util.FormatPrState(pr.State), pr.ID, pr.Title, pr.Source.Branch.Name, pr.Destination.Branch.Name, pr.Author.Nickname)
			if status {
				fmt.Printf(" %s", util.FormatPipelineState(pr.Status.State))
			}
			if participants {
				var outputStr = []string{}
				for _, participant := range pr.Participants {
					if participant.Approved {
						outputStr = append(outputStr, fmt.Sprintf("\033[1;32m✓ %s\033[m", participant.User.DisplayName))
					} else {
						outputStr = append(outputStr, fmt.Sprintf("\033[0;37m%s\033[m", participant.User.DisplayName))
					}
				}
				fmt.Printf("\n       \033[37mComments: %d\033[m ( %s )", pr.CommentCount, strings.Join(outputStr, ", "))
			}
			fmt.Println()
			count++
		}
		if count == 0 {
			fmt.Printf("\n  No pull requests for \033[1;36m%s\033[m\n\n", viper.GetString("repo"))
		}
	},
}

func init() {
	// filter
	ListCmd.Flags().StringP("author", "a", "", "filter by author nick name (full nickname is needed due to an API limitation from bitbucket)")
	ListCmd.Flags().String("search", "", "search pull request with query")
	ListCmd.Flags().StringArrayP("state", "s", []string{string(api.OPEN)}, `filter by state. Default: "open". Multiple of these options can be given
	possible options: "open", "merged", "declined" or "superseded"`)
	ListCmd.RegisterFlagCompletionFunc("state", stateCompletion)
	ListCmd.Flags().String("destination", "", "filter by destination branch.")
	ListCmd.RegisterFlagCompletionFunc("destination", branchCompletion)
	ListCmd.Flags().String("source", "", "filter by source branch.")
	ListCmd.RegisterFlagCompletionFunc("source", branchCompletion)
	ListCmd.Flags().Bool("all", false, "return pull request with all possible states.")

	ListCmd.Flags().Int("pages", 1, "number of pages with results to retrieve")
	ListCmd.Flags().BoolP("status", "S", false, "include status of each pull request on the result. (the result will be slower)")
	ListCmd.Flags().BoolP("participants", "p", false, "include participant and comment data for each pull request on the result. (the result will be slower)")
}

func stateCompletion(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"open", "merged", "declined", "superseded"}, cobra.ShellCompDirectiveDefault
}

func branchCompletion(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return util.ListBranches(), cobra.ShellCompDirectiveDefault
}
