package pr

import (
	"bb/api"
	"bb/util"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

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

		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("body")
		source, _ := cmd.Flags().GetString("source")
		destination, _ := cmd.Flags().GetString("destination")
		close_source, _ := cmd.Flags().GetBool("close_source")

		// if no options given ask for what to change
		var existingPr api.PullRequest
		if !cmd.Flags().Changed("title") &&
			!cmd.Flags().Changed("body") {
			existingPr = <-api.GetPr(repo, id)
			title, description = readTitleAndDescription(existingPr)
		}

		newpr := api.CreatePullRequest{
			Title:       title,
			Description: description,
			CloseSource: close_source,
		}
		// TODO handle these empty fields
		newpr.Destination.Branch.Name = source
		newpr.Destination.Branch.Name = destination
		newpr.Reviewers = nil

		pr := api.UpdatePr(repo, id, newpr)

		fmt.Printf("\n%s \033[1;32m#%d\033[m \033[1;37m%s\033[m\n", util.FormatPrState(pr.State), pr.ID, pr.Title)
		fmt.Printf("\033[37m  opened by %s, %d comments, last updated: %s\033[m\n\n", pr.Author.Nickname, pr.CommentCount, util.TimeAgo(pr.UpdatedOn))
		if pr.Description != "" {
			fmt.Printf("%s\n\n", pr.Description)
		}
	},
}

func init() {
	EditCmd.Flags().StringP("title", "t", "", "title for the pull request. \033[31mNot implemented\033[m")
	EditCmd.Flags().StringP("body", "b", "", "description for the pull request. \033[31mNot implemented\033[m")
	EditCmd.Flags().StringP("source", "s", "", "source branch. Defaults to current branch. \033[31mNot implemented\033[m")
	EditCmd.Flags().StringP("destination", "d", "dev", "description for the pull request: Defaults to dev. \033[31mNot implemented\033[m")
	EditCmd.Flags().BoolP("close-source", "c", false, "close source branch. \033[31mNot implemented\033[m")
	EditCmd.Flags().StringArrayP("reviewer", "r", []string{}, "add reviewer by their name. \033[31mNot implemented\033[m")
}

func readTitleAndDescription(pr api.PullRequest) (string, string) {
	tmpFile, err := ioutil.TempFile("/tmp", "bitbucket-pr-edit-")
	cobra.CheckErr(err)

	tmpFile.WriteString(pr.Title + "\n\n")
	// TODO maybe add a delimiter ?
	tmpFile.WriteString(pr.Description)
	tmpFile.Seek(0, 0)

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
	fullFile, err := ioutil.ReadAll(tmpFile)

	title := ""
	description := ""
	lines := strings.Split(string(fullFile), "\n")
	if len(lines) > 0 {
		title = lines[0]
	}
	if len(lines) > 2 {
		description = strings.Join(lines[2:], "\n")
	}
	return title, description
}
