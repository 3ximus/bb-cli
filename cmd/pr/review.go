package pr

import (
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

		fmt.Println(id, repo)
	},
}

func init() {
	ReviewCmd.Flags().BoolP("merge", "m", false, "Merge pull request. \033[31mNot implemented\033[m")
	ReviewCmd.Flags().BoolP("approve", "a", false, "Approve pull request. \033[31mNot implemented\033[m")
	ReviewCmd.Flags().BoolP("unnaprove", "u", false, "Unnaprove pull request. \033[31mNot implemented\033[m")
	ReviewCmd.Flags().BoolP("decline", "d", false, "Decline pull request. \033[31mNot implemented\033[m")
	ReviewCmd.Flags().BoolP("request-changes", "c", false, "Request changes to the pull request. \033[31mNot implemented\033[m")
	ReviewCmd.Flags().BoolP("unrequest-changes", "U", false, "Remove request changes status from pull request. \033[31mNot implemented\033[m")
	ReviewCmd.MarkFlagsMutuallyExclusive("merge", "approve", "unnaprove", "decline", "request-changes", "unrequest-changes")
}
