package pr

import (
	"bb/api"
	"bb/util"
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a pull request on a repository",
	Args:  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlag("include_branch_name", cmd.Flags().Lookup("include-branch-name"))
		cobra.CheckErr(err)
	},
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")
		scanner := bufio.NewScanner(os.Stdin)

		// set account id if it doesn't exist
		authorId := viper.GetString("account_id")
		if authorId == "" {
			// TODO make this into an async call that we can retrieve the result later
			user := api.GetUser()
			viper.Set("account_id", user.AccountId)
			// TODO Don't do this because it permanently saves the value from "repo"
			// and subsequent calls will only use that value
			// viper.WriteConfig()
			authorId = user.AccountId
		}

		// load reviewers
		membersChannel := api.GetWorkspaceMembers(strings.Split(repo, "/")[0])
		reviewersChannel := api.GetReviewers(repo)

		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("body")
		if title == "" {
			fmt.Print("? \033[1;35mTitle \033[m")
			scanner.Scan()
			title = scanner.Text()
			if description == "" {
				description = readDescription(scanner)
			}
		}
		source, _ := cmd.Flags().GetString("source")
		destination, _ := cmd.Flags().GetString("destination")
		close_source, _ := cmd.Flags().GetBool("close_source")
		include_branch_name := viper.GetBool("include_branch_name")

		if source == "" {
			var err error
			source, err = util.GetCurrentBranch()
			cobra.CheckErr(err)
		}

		// select reviewers
		members, reviewers := <-membersChannel, <-reviewersChannel
		if len(reviewers) == 0 {
			reviewers = members // use members instead of reviewers
		}
		// TODO filter author id out
		reviewersIndexes := chooseReviewers(reviewers)

		if include_branch_name {
			re := regexp.MustCompile(api.JiraIssueKeyRegex)
			key := re.FindString(source)
			title = key + " " + title
		}

		// create dto
		newpr := api.CreatePullRequest{
			Title:       title,
			Description: description,
			CloseSource: close_source,
		}
		newpr.Source = &api.Branch{}
		newpr.Source.Branch.Name = source
		newpr.Destination = &api.Branch{}
		newpr.Destination.Branch.Name = destination
		for _, idx := range reviewersIndexes {
			newpr.Reviewers = append(newpr.Reviewers, struct {
				AccountId string `json:"account_id"`
			}{AccountId: members[idx].AccountId})
		}
		if len(reviewersIndexes) == 0 {
			newpr.Reviewers = []struct {
				AccountId string `json:"account_id"`
			}{}
		}

		// confirm pr
		fmt.Printf("\033[1;37m%s\033[m  \033[1;34m[ %s → %s ]\033[m\n", newpr.Title, newpr.Source.Branch.Name, newpr.Destination.Branch.Name)
		if newpr.Description != "" {
			fmt.Printf("%s\n", newpr.Description)
		}
		if len(reviewersIndexes) > 0 {
			fmt.Println("Reviewers:")
			for _, idx := range reviewersIndexes {
				fmt.Printf("  - %s \033[37m( ID: %s )\033[m\n", members[idx].DisplayName, members[idx].AccountId)
			}
		}
		fmt.Print("? \033[1;35mCreate this PR ? [y/n]\033[m ")
		scanner.Scan()
		if strings.TrimSpace(strings.ToLower(scanner.Text())) != "y" {
			return
		}

		// send create request
		pr := api.PostPr(repo, newpr)

		fmt.Printf("\n%s \033[1;32m#%d\033[m \033[1;37m%s\033[m\n", util.FormatPrState(pr.State), pr.ID, pr.Title)
		fmt.Printf("\033[37m  opened by %s, %d comments, last updated: %s\033[m\n\n", pr.Author.Nickname, pr.CommentCount, util.TimeAgo(pr.UpdatedOn))
		if pr.Description != "" {
			fmt.Printf("%s\n\n", pr.Description)
		}
	},
}

func init() {
	CreateCmd.Flags().StringP("title", "t", "", "title for the pull request")
	CreateCmd.Flags().StringP("body", "b", "", "description for the pull request")
	CreateCmd.Flags().StringP("source", "s", "", "source branch. Defaults to current branch")
	CreateCmd.Flags().StringP("destination", "d", "dev", "description for the pull request: Defaults to dev")
	CreateCmd.RegisterFlagCompletionFunc("source", branchCompletion)
	CreateCmd.RegisterFlagCompletionFunc("destination", branchCompletion)
	CreateCmd.Flags().BoolP("close-source", "c", true, "close source branch")
	CreateCmd.Flags().StringArrayP("reviewer", "r", []string{}, "add reviewer by their name. \033[31mNot implemented\033[m")
	CreateCmd.Flags().BoolP("include-branch-name", "i", false, "include branch name in the pull request name")
}

func readDescription(scanner *bufio.Scanner) string {
	fmt.Print("? \033[1;35mAdd body ? [y/n]\033[m ")
	scanner.Scan()
	if strings.TrimSpace(strings.ToLower(scanner.Text())) == "y" {
		tmpFile, err := os.CreateTemp("/tmp", "bitbucket-pr-body-")
		cobra.CheckErr(err)
		defer os.Remove(tmpFile.Name())
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		cmd := exec.Command(editor, tmpFile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Start()
		cobra.CheckErr(err)
		err = cmd.Wait()
		cobra.CheckErr(err)
		description, err := io.ReadAll(tmpFile)
		return string(description)
	}
	fmt.Println()
	return ""
}

func chooseReviewers(reviewers []api.User) []int {
	return util.SelectFZF(reviewers, "Reviewers > ", func(i int) string {
		return fmt.Sprintf("%s", reviewers[i].Nickname)
	})
}
