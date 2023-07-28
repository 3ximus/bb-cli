package pr

import (
	"bb/api"
	"bb/util"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ReviewCmd = &cobra.Command{
	Use:   "review [ID]",
	Short: "Review a pull request (merge, approve, unnaprove, decline ...)",
	Long: `Merge, approve, unnaprove, decline or request/unrequest changes in a pull request
	If no ID is given the operation will be applied to the first PR found for the current branch`,
	ValidArgsFunction: func(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var opt = []string{}
		for pr := range api.GetPrList(util.GetCurrentRepo(), []string{string(api.OPEN)}, "", "", "", "", 1, false, false) {
			opt = append(opt, fmt.Sprint(pr.ID))
		}
		return opt, cobra.ShellCompDirectiveDefault
	},
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")

		// TODO allow to give source branch name instead of just using current branch
		var id int
		var err error
		if len(args) == 0 {
			branch, err := util.GetCurrentBranch()
			cobra.CheckErr(err)
			// retrieve id of pr for current branch
			pr := <-api.GetPrList(repo, []string{string(api.OPEN), string(api.MERGED), string(api.DECLINED), string(api.SUPERSEDED)}, "", "", branch, "", 1, false, false)
			if pr.ID == 0 {
				cobra.CheckErr("No pr found for this branch")
			}
			id = pr.ID // get the first one's ID
		} else {
			id, err = strconv.Atoi(args[0])
			cobra.CheckErr(err)
		}

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
		merge, _ := cmd.Flags().GetBool("merge")
		if merge {
			message,_ := cmd.Flags().GetString("message")
			api.MergePr(repo, id, message)
			fmt.Printf("\033[1;35mMerge\033[m pull request #%d\n", id)
		}
		requestChanges, _ := cmd.Flags().GetBool("request-changes")
		if requestChanges {
			api.RequestChangesPr(repo, id)
			fmt.Printf("\033[1;34mRequested changes\033[m for pull request #%d\n", id)
		}
		unrequestChanges, _ := cmd.Flags().GetBool("unrequest-changes")
		if unrequestChanges {
			fmt.Printf("Not implemented")
		}

		if !merge && !approve && !unnaprove && !decline && !requestChanges && !unrequestChanges {
			fmt.Println("No operation selected")
			cmd.Help()
		}
	},
}

func init() {
	ReviewCmd.Flags().BoolP("merge", "m", false, "Merge pull request")
	ReviewCmd.Flags().BoolP("approve", "a", false, "Approve pull request")
	ReviewCmd.Flags().BoolP("unnaprove", "u", false, "Unnaprove pull request")
	ReviewCmd.Flags().BoolP("decline", "d", false, "Decline pull request")
	ReviewCmd.Flags().BoolP("request-changes", "c", false, "Request changes to the pull request")
	ReviewCmd.Flags().BoolP("unrequest-changes", "U", false, "Remove request changes status from pull request. \033[31mNot implemented\033[m")
	ReviewCmd.MarkFlagsMutuallyExclusive("merge", "approve", "unnaprove", "decline", "request-changes", "unrequest-changes")

	ReviewCmd.Flags().String("message", "", "Attach message to action (merge)")
}
