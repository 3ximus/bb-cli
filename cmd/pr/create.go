package pr

import (
	"bb/api"
	"bufio"
	"fmt"
	"github.com/ldez/go-git-cmd-wrapper/v2/branch"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a pull request on a repository",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")

		scanner := bufio.NewScanner(os.Stdin)

		title, _ := cmd.Flags().GetString("title")
		if title == "" {
			fmt.Print("\033[1;37m? Title \033[m")
			scanner.Scan()
			title = scanner.Text()
		}
		source, _ := cmd.Flags().GetString("source")
		destination, _ := cmd.Flags().GetString("destination")
		description, _ := cmd.Flags().GetString("body")
		if description == "" {
			// TODO launch text editor instead
			fmt.Print("\033[1;37m? Body \033[m")
			scanner.Scan()
			description = scanner.Text()
		}
		close_source, _ := cmd.Flags().GetBool("close_source")
		api.PostPr(repo, source, destination, title, description, close_source)
	},
}

func init() {
	CreateCmd.Flags().StringP("title", "t", "", "title for the pull request")
	CreateCmd.Flags().StringP("body", "b", "", "description for the pull request")
	CreateCmd.Flags().StringP("source", "s", getCurrentBranch(), "source branch. Defaults to current branch")
	CreateCmd.Flags().StringP("destination", "d", "dev", "description for the pull request: Defaults to dev")
	CreateCmd.Flags().BoolP("close-source", "c", true, "close source branch")
}

func getCurrentBranch() string {
	branch, err := git.Branch(branch.ShowCurrent)
	cobra.CheckErr(err)
	return strings.Trim(branch, "\n")
}
