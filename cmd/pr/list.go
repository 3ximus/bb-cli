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
	Use:     "list",
	Short:   "List pull requests from a repository",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
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
		target, _ := cmd.Flags().GetString("target")
		status, _ := cmd.Flags().GetBool("status")
		participants, _ := cmd.Flags().GetBool("participants")

		count := 0
		for pr := range api.GetPrList(viper.GetString("repo"), states, author, search, source, target, pages, status, participants) {
			util.Printf("%s \033[1;32m#%d\033[m %s \033[1;34m[ %s \033[m→\033[1;34m %s ]\033[m \033[33m%s\033[m", util.FormatPrState(pr.State), pr.ID, pr.Title, pr.Source.Branch.Name, pr.Destination.Branch.Name, pr.Author.Nickname)
			if status {
				util.Printf(" %s", util.FormatPipelineStatus(pr.Status.State))
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
				util.Printf("\n       \033[37mComments: %d\033[m ( %s )", pr.CommentCount, strings.Join(outputStr, ", "))
			}
			fmt.Println()
			count++
		}
		if count == 0 {
			util.Printf("No pull requests for \033[1;36m%s\033[m\n", viper.GetString("repo"))
		}
	},
}

func init() {
	ListCmd.Flags().StringP("author", "a", "", "filter by author nick name (full nickname is needed due to an API limitation from bitbucket)")
	ListCmd.Flags().String("search", "", "search pull request with query")
	ListCmd.Flags().StringArrayP("state", "s", []string{string(api.OPEN)}, `filter by state. Default: "open". Multiple of these options can be given
	possible options: "open", "merged", "declined" or "superseded"`)
	ListCmd.RegisterFlagCompletionFunc("state", stateCompletion)
	ListCmd.Flags().String("target", "", "filter by target branch.")
	ListCmd.RegisterFlagCompletionFunc("target", util.BranchCompletion)
	ListCmd.Flags().String("source", "", "filter by source branch.")
	ListCmd.RegisterFlagCompletionFunc("source", util.BranchCompletion)
	ListCmd.Flags().Bool("all", false, "return pull request with all possible states.")

	ListCmd.Flags().Int("pages", 1, "number of pages with results to retrieve")
	ListCmd.Flags().BoolP("status", "S", false, "include status of each pull request on the result. (the result will be slower)")
	ListCmd.Flags().BoolP("participants", "p", false, "include participant and comment data for each pull request on the result. (the result will be slower)")
}

func stateCompletion(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"open", "merged", "declined", "superseded"}, cobra.ShellCompDirectiveDefault
}
