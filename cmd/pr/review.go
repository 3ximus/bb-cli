package pr

import (
	"bb/api"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ReviewCmd = &cobra.Command{
	Use:   "review ID",
	Short: "Review a pull request",
	Long:  "Merge, approve, unnaprove, decline or request/unrequest changes in a pull request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		cobra.CheckErr(err)

		repo := viper.GetString("repo")

		approve, _ := cmd.Flags().GetBool("approve")
		if approve {
			api.ApprovePr(repo, id)
			fmt.Printf("Pull request #%d \033[1;32mApproved\033[m\n", id)
		}
		unnaprove, _ := cmd.Flags().GetBool("unnaprove")
		if unnaprove {
			api.UnnaprovePr(repo, id)
			fmt.Printf("Pull request #%d \033[1;33mUnnaproved\033[m\n", id)
		}
		decline, _ := cmd.Flags().GetBool("decline")
		if decline {
			api.DeclinePr(repo, id)
			fmt.Printf("Pull request #%d \033[1;31mDeclined\033[m\n", id)
		}
		requestChanges, _ := cmd.Flags().GetBool("request-changes")
		if requestChanges {
			api.RequestChangesPr(repo, id)
			fmt.Printf("\033[1;34mRequested changes\033[m for pull request #%d\n", id)
		}
	},
}

func init() {
	ReviewCmd.Flags().BoolP("merge", "m", false, "Merge pull request. \033[31mNot implemented\033[m")
	ReviewCmd.Flags().BoolP("approve", "a", false, "Approve pull request")
	ReviewCmd.Flags().BoolP("unnaprove", "u", false, "Unnaprove pull request")
	ReviewCmd.Flags().BoolP("decline", "d", false, "Decline pull request")
	ReviewCmd.Flags().BoolP("request-changes", "c", false, "Request changes to the pull request")
	ReviewCmd.Flags().BoolP("unrequest-changes", "U", false, "Remove request changes status from pull request. \033[31mNot implemented\033[m")
	ReviewCmd.MarkFlagsMutuallyExclusive("merge", "approve", "unnaprove", "decline", "request-changes", "unrequest-changes")
}
