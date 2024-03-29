package pr

import (
	"bb/api"
	"bb/util"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var EditCmd = &cobra.Command{
	Use:   "edit ID",
	Short: "Edit details of a pull request",
	Long: `Allows edits to an existing pull request
	If no options are given to edit title or description it will open your EDITOR to write any changes to them.
	By default title is on first line and description on the lines bellow`,
	Args: cobra.MaximumNArgs(1),
	ValidArgsFunction: func(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var opt = []string{}
		for pr := range api.GetPrList(util.GetCurrentRepo(), []string{string(api.OPEN)}, "", "", "", "", 1, false, false) {
			opt = append(opt, fmt.Sprint(pr.ID))
		}
		return opt, cobra.ShellCompDirectiveDefault
	},
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")

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

		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("body")
		source, _ := cmd.Flags().GetString("source")
		target, _ := cmd.Flags().GetString("target")
		close_source, _ := cmd.Flags().GetBool("close_source")

		// if no options given ask for what to change
		existingPr := <-api.GetPr(repo, id)
		if !cmd.Flags().Changed("title") && !cmd.Flags().Changed("body") {
			title, description = readTitleAndDescription(existingPr)
		}
		if title == "" && !cmd.Flags().Changed("title") {
			title = existingPr.Title
		}
		if description == "" && !cmd.Flags().Changed("description") {
			description = existingPr.Description
		}

		if !cmd.Flags().Changed("close_source") {
			close_source = existingPr.CloseSource
		}

		newpr := api.CreatePullRequestBody{
			Title:       title,
			Description: description,
			CloseSource: close_source,
		}
		if source == "" {
			newpr.Source = nil
		} else {
			newpr.Source = &api.Branch{}
			newpr.Source.Branch.Name = source
		}
		if target == "" {
			newpr.Destination = nil
		} else {
			newpr.Destination = &api.Branch{}
			newpr.Destination.Branch.Name = target
		}
		newpr.Reviewers = nil

		pr := api.UpdatePr(repo, id, newpr)

		fmt.Printf("\n%s \033[1;32m#%d\033[m \033[1;37m%s\033[m  \033[1;34m[ %s → %s]\033[m\n", util.FormatPrState(pr.State), pr.ID, pr.Title, pr.Source.Branch.Name, pr.Destination.Branch.Name)
		fmt.Printf("\033[37m  opened by %s, %d comments, last updated: %s\033[m\n\n", pr.Author.Nickname, pr.CommentCount, util.TimeAgo(pr.UpdatedOn))
		if pr.Description != "" {
			fmt.Printf("%s\n\n", pr.Description)
		}
	},
}

func init() {
	EditCmd.Flags().StringP("title", "t", "", "title for the pull request")
	EditCmd.Flags().StringP("body", "b", "", "description for the pull request")
	EditCmd.Flags().String("source", "", "source branch. Defaults to current branch")
	EditCmd.Flags().String("target", "", "description for the pull request")
	EditCmd.RegisterFlagCompletionFunc("source", util.BranchCompletion)
	EditCmd.RegisterFlagCompletionFunc("target", util.BranchCompletion)
	EditCmd.Flags().BoolP("close-source", "c", false, "close source branch")
	EditCmd.Flags().StringArrayP("reviewer", "r", []string{}, "add reviewer by their name. \033[31mNot implemented\033[m")
}

func readTitleAndDescription(pr api.PullRequest) (string, string) {
	tmpFile, err := os.CreateTemp("/tmp", "bitbucket-pr-edit-")
	cobra.CheckErr(err)

	tmpFile.WriteString(pr.Title + "\n\n")
	// TODO maybe add a delimiter ?
	tmpFile.WriteString(pr.Description)
	tmpFile.Seek(0, 0)

	defer os.Remove(tmpFile.Name())
	util.OpenInEditor(tmpFile)
	fullFile, err := io.ReadAll(tmpFile)

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
