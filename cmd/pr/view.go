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
	Use:   "view ID",
	Short: "View a pull request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		cobra.CheckErr(err)

		repo := viper.GetString("repo")
		prChannel := api.GetPr(repo, id)
		statusesChannel := api.GetPrStatuses(repo, id)

		pr := <-prChannel
		fmt.Printf("\n%s \033[1;32m#%d\033[m \033[1;37m%s\033[m\n", util.FormatPrState(pr.State), pr.ID, pr.Title)
		fmt.Printf("\033[37m  opened by %s, %d comments, last updated: %s\033[m\n\n", pr.Author.Nickname, pr.CommentCount, util.TimeAgo(pr.UpdatedOn))
		if pr.Description != "" {
			fmt.Printf("%s\n\n", pr.Description)
		}
		fmt.Println("Pipelines:")
		for _, pipeline := range <-statusesChannel {
			fmt.Printf("%s %s \033[37m(%s)\033[m\n", util.FormatPipelineState(pipeline.State), pipeline.Name, pipeline.RefName)
			fmt.Printf("  \033[37m%s\033[m\n", pipeline.Url)
		}
		fmt.Println()
	},
}

func init() {
	ViewCmd.Flags().BoolP("comments", "c", false, "View comments. \033[31mNot implemented\033[m")
	ViewCmd.Flags().BoolP("commits", "C", false, "View commits. \033[31mNot implemented\033[m")
	ViewCmd.Flags().BoolP("diff", "d", false, "View diff. \033[31mNot implemented\033[m")
}
