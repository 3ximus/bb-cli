package pr

import (
	"bb/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a pull request on a repository",
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")
		source, _ := cmd.Flags().GetString("source")
		destination, _ := cmd.Flags().GetString("destination")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("body")
		close_source, _ := cmd.Flags().GetBool("close_source")
		api.PostPr(repo, source, destination, title, description, close_source)
	},
}

func init() {
	CreateCmd.Flags().StringP("title", "t", "", "Title for the pull request")
	CreateCmd.Flags().StringP("body", "b", "", "Description for the pull request")
	CreateCmd.Flags().StringP("source", "s", "", "Source branch: Defaults to current branch")
	CreateCmd.Flags().StringP("destination", "d", "dev", "Description for the pull request: Defaults to dev")
	CreateCmd.Flags().BoolP("close-source", "c", true, "Close source branch")
}
