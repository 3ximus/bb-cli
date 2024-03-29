package pr

import (
	"bb/api"
	"bb/util"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ViewCmd = &cobra.Command{
	Use:   "view [ID]",
	Short: "View details of a pull request",
	Long: `View details of a pull request from given ID.
	If no ID is given we'll try to find an open pull request that has it's source as the current branch`,
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
		showComments, _ := cmd.Flags().GetBool("comments")
		targetBranch, _ := cmd.Flags().GetString("target")
		sourceBranch, _ := cmd.Flags().GetString("source")

		var id int
		var err error
		if len(args) == 0 {
			var err error
			if targetBranch == "" && sourceBranch == "" {
				sourceBranch, err = util.GetCurrentBranch()
				targetBranch = ""
				cobra.CheckErr(err)
			}
			// retrieve id of pr for current branch
			pr := <-api.GetPrList(repo, []string{string(api.OPEN), string(api.MERGED), string(api.DECLINED), string(api.SUPERSEDED)}, "", "", sourceBranch, targetBranch, 1, false, false)
			if pr.ID == 0 {
				cobra.CheckErr(fmt.Sprintf("No pull request found for branches (source: '%s', target: '%s')", sourceBranch, targetBranch))
			}
			id = pr.ID // get the first one's ID
		} else {
			id, err = strconv.Atoi(args[0])
			cobra.CheckErr(err)
		}

		statusesChannel := api.GetPrStatuses(repo, id)
		commentsChannel := api.GetPrComments(repo, id)

		// BASIC INFO

		pr := <-api.GetPr(repo, id)
		util.Printf("\n%s \033[1;32m#%d\033[m \033[1;37m%s\033[m  \033[1;34m[ %s → %s]\033[m\n", util.FormatPrState(pr.State), pr.ID, pr.Title, pr.Source.Branch.Name, pr.Destination.Branch.Name)
		util.Printf("\033[37m  opened by %s, %d comments, last updated: %s\033[m\n", pr.Author.Nickname, pr.CommentCount, util.TimeAgo(pr.UpdatedOn))
		util.Printf("\033[37m  reviewers: \n")
		for _, participant := range pr.Participants {
			if participant.Approved {
				util.Printf("    \033[1;32m✓ %s\n", participant.User.DisplayName)
			} else {
				util.Printf("    \033[0;37m%s\n", participant.User.DisplayName)
			}
		}
		util.Printf("\033[m\n")
		if pr.Description != "" {
			util.Printf("%s\n\n", pr.Description)
		}

		web, _ := cmd.Flags().GetBool("web")
		if web {
			util.OpenInBrowser(pr.Links.Html.Href)
			return
		}

		// PIPELINES

		pipelines := <-statusesChannel
		if pipelines != nil && len(pipelines) > 0 {
			fmt.Println("Pipelines:")
			for _, pipeline := range pipelines {
				util.Printf("%s %s \033[37m(%s)\033[m\n", util.FormatPipelineStatus(pipeline.State), pipeline.Name, pipeline.RefName)
				util.Printf("  \033[37m%s\033[m\n", pipeline.Url)
			}
			fmt.Println()
		}

		if showComments {
			comments := <-commentsChannel
			for _, comment := range comments {
				util.Printf("%s %s", comment.Content.Raw, comment.User.DisplayName)
			}
		}

	},
}

func init() {
	ViewCmd.Flags().String("target", "", "filter by target branch.")
	ViewCmd.Flags().String("source", "", "filter by source branch.")
	ViewCmd.RegisterFlagCompletionFunc("target", util.BranchCompletion)
	ViewCmd.RegisterFlagCompletionFunc("source", util.BranchCompletion)
	ViewCmd.Flags().BoolP("comments", "c", false, "View comments")
	ViewCmd.Flags().Bool("web", false, "Open in the browser.")
}
