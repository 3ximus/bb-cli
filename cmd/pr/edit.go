package pr

import (
	"bb/util"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var EditCmd = &cobra.Command{
	Use:   "edit ID",
	Short: "Edit details of a pull request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		cobra.CheckErr(err)

		repo := viper.GetString("repo")

		fmt.Println(id, repo)
	},
}

func init() {
	EditCmd.Flags().StringP("title", "t", "", "title for the pull request")
	EditCmd.Flags().StringP("body", "b", "", "description for the pull request")
	EditCmd.Flags().StringP("source", "s", util.GetCurrentBranch(), "source branch. Defaults to current branch")
	EditCmd.Flags().StringP("destination", "d", "dev", "description for the pull request: Defaults to dev")
	EditCmd.Flags().BoolP("close-source", "c", true, "close source branch")
	EditCmd.Flags().StringArrayP("reviewer", "r", []string{}, "add reviewer by their name. \033[31mNot implemented\033[m")
}
