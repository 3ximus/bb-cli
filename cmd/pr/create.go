package pr

import (
	"bb/api"
	"bb/util"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	// "github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a pull request on a repository",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")
		scanner := bufio.NewScanner(os.Stdin)

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

		// select reviewers
		members, reviewers := <-membersChannel, <-reviewersChannel
		if len(reviewers) > 0 {
			members = append(reviewers)
		}
		reviewersIndexes := chooseReviewers(members)

		// create dto
		newpr := api.CreatePullRequest{
			Title:       title,
			Description: description,
			CloseSource: close_source,
		}
		newpr.Source.Branch.Name = source
		newpr.Destination.Branch.Name = destination
		for _, idx := range reviewersIndexes {
			newpr.Reviewers = append(newpr.Reviewers, struct {
				AccountId string `json:"account_id"`
			}{AccountId: members[idx].AccountId})
		}

		// confirm pr
		fmt.Printf("\033[1;37m%s\033[m  \033[1;34m[ %s → %s]\033[m\n", newpr.Title, newpr.Source.Branch.Name, newpr.Destination.Branch.Name)
		fmt.Println("Reviewers:")
		for _, idx := range reviewersIndexes {
			fmt.Printf("  - %s \033[37m( ID: %s )\033[m\n", members[idx].DisplayName, members[idx].AccountId)
		}
		if newpr.Description != "" {
			fmt.Printf("%s\n", newpr.Description)
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
	CreateCmd.Flags().StringP("source", "s", util.GetCurrentBranch(), "source branch. Defaults to current branch")
	CreateCmd.Flags().StringP("destination", "d", "dev", "description for the pull request: Defaults to dev")
	CreateCmd.Flags().BoolP("close-source", "c", true, "close source branch")
}

func readDescription(scanner *bufio.Scanner) string {
	fmt.Print("? \033[1;35mAdd body ? [y/n]\033[m ")
	scanner.Scan()
	if strings.TrimSpace(strings.ToLower(scanner.Text())) == "y" {
		tmpFile, err := ioutil.TempFile("/tmp", "bitbucket-pr-body-")
		cobra.CheckErr(err)
		defer os.Remove(tmpFile.Name())
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim"
		}
		cmd := exec.Command(editor, tmpFile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Start()
		cobra.CheckErr(err)
		err = cmd.Wait()
		cobra.CheckErr(err)
		description, err := ioutil.ReadAll(tmpFile)
		return string(description)
	}
	fmt.Println()
	return ""
}

func chooseReviewers(reviewers []api.User) []int {
	if len(reviewers) == 0 {
		return []int{}
	}

	return UseExternalFZF(reviewers, func(i int) string {
		return fmt.Sprintf("%s", reviewers[i].Nickname)
	})

	// This would be good to not depend on external fzf but its ugly... Maybe just use it as backup ?
	// indexes, err := fuzzyfinder.FindMulti(reviewers, func(i int) string {
	// 	return fmt.Sprintf("%s (%s)", reviewers[i].Nickname, reviewers[i].DisplayName)
	// }, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop))
	// cobra.CheckErr(err)
	// return indexes
}

func UseExternalFZF(list []api.User, toString func(int) string) []int {
	input := ""
	for i := range list {
		input += fmt.Sprintf("%d %s\n", i, toString(i))
	}
	cmd := exec.Command("fzf", "-m", "--height", "20%", "--reverse", "--with-nth", "2..")
	var selectionBuffer strings.Builder
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &selectionBuffer
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	cobra.CheckErr(err)
	err = cmd.Wait()
	cobra.CheckErr(err)

	var result []int
	for _, r := range strings.Split(selectionBuffer.String(), "\n") {
		if r == "" {
			continue
		}
		idx, err := strconv.Atoi(strings.Split(r, " ")[0])
		cobra.CheckErr(err)
		result = append(result, idx)
	}
	return result
}
