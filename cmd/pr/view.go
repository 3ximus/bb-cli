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

		var id int
		var err error
		if len(args) == 0 {
			branch := util.GetCurrentBranch()
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

		statusesChannel := api.GetPrStatuses(repo, id)

		// BASIC INFO

		pr := <-api.GetPr(repo, id)
		fmt.Printf("\n%s \033[1;32m#%d\033[m \033[1;37m%s\033[m  \033[1;34m[ %s → %s]\033[m\n", util.FormatPrState(pr.State), pr.ID, pr.Title, pr.Source.Branch.Name, pr.Destination.Branch.Name)
		fmt.Printf("\033[37m  opened by %s, %d comments, last updated: %s\033[m\n", pr.Author.Nickname, pr.CommentCount, util.TimeAgo(pr.UpdatedOn))
		fmt.Print("\033[37m  reviewers: \n")
		for _, participant := range pr.Participants {
			if participant.Approved {
				fmt.Printf("    \033[1;32m✓ %s\n", participant.User.DisplayName)
			} else {
				fmt.Printf("    \033[0;37m%s\n", participant.User.DisplayName)
			}
		}
		fmt.Print("\033[m\n")
		if pr.Description != "" {
			fmt.Printf("%s\n\n", pr.Description)
		}

		web, _ := cmd.Flags().GetBool("web")
		if web {
			util.OpenInBrowser(pr.Links.Html.Href)
			return
		}

		// PIPELINES

		pipelines := <-statusesChannel
		if len(pipelines) != 0 {
			fmt.Println("Pipelines:")
			for _, pipeline := range pipelines {
				fmt.Printf("%s %s \033[37m(%s)\033[m\n", util.FormatPipelineState(pipeline.State), pipeline.Name, pipeline.RefName)
				fmt.Printf("  \033[37m%s\033[m\n", pipeline.Url)
			}
			fmt.Println()
		}

	},
}

func init() {
	ViewCmd.Flags().BoolP("comments", "c", false, "View comments. \033[31mNot implemented\033[m")
	ViewCmd.Flags().BoolP("commits", "C", false, "View commits. \033[31mNot implemented\033[m")
	ViewCmd.Flags().BoolP("diff", "d", false, "View diff. \033[31mNot implemented\033[m")
	ViewCmd.Flags().Bool("web", false, "Open in the browser.")
}
