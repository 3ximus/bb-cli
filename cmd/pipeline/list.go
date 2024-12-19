package pipeline

import (
	"bb/api"
	"bb/util"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pipelines from a repository",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		nResults, _ := cmd.Flags().GetInt("number-results")
		showAuthor, _ := cmd.Flags().GetBool("author")
		targetBranch, _ := cmd.Flags().GetString("target")

		useFZF, _ := cmd.Flags().GetBool("fzf")
		useFZFInternal, _ := cmd.Flags().GetBool("fzf-internal")
		if useFZF {
			util.ReplaceListWithFzf("--ansi --reverse --height 40% --read0 --prompt 'View > ' --info-command 'echo -e $FZF_POS/$FZF_INFO' --info inline --bind 'enter:become(" + os.Args[0] + " pipeline view -R "+viper.GetString("repo")+" {2})' --bind 'ctrl-w:become(" + os.Args[0] + " pipeline view -R "+viper.GetString("repo")+" --web {2})'")
			return
		}

		for pipeline := range api.GetPipelineList(viper.GetString("repo"), nResults, targetBranch) {
			if pipeline.State.Result.Name == "" {
				util.Printf("%s", util.FormatPipelineStatus(pipeline.State.Name))
			} else {
				util.Printf("%s", util.FormatPipelineStatus(pipeline.State.Result.Name))
			}
			util.Printf(" \033[1;32m#%d\033[m ", pipeline.BuildNumber)
			if pipeline.Target.Source != "" {
				util.Printf("%s \033[1;34m[ %s → %s ]\033[m", pipeline.Target.PullRequest.Title, pipeline.Target.Source, pipeline.Target.Destination)
			} else {
				util.Printf("\033[1;34m[ %s ]\033[m", pipeline.Target.RefName)
			}

			util.Printf(" \033[37m%s (%s)\033[m", util.TimeDuration(time.Duration(pipeline.DurationInSeconds*1000000000)), util.TimeAgo(pipeline.CreatedOn))

			if showAuthor {
				util.Printf("\n        \033[33m%s\033[m \033[37mTrigger: %s\033[m", pipeline.Author.DisplayName, pipeline.Trigger.Name) //  \033[37mComments: %d\033[m",
			}

			endChar := "\n"
			if useFZFInternal {
				endChar = "\x00"
			}
			util.Printf(endChar)
		}
	},
}

func init() {
	ListCmd.Flags().IntP("number-results", "n", 10, "max number of results retrieve (max: 100)")
	ListCmd.Flags().BoolP("author", "a", false, "show author information")
	ListCmd.Flags().String("target", "", "filter target branch of pipeline")
	if util.CommandExists("fzf") {
		ListCmd.Flags().Bool("fzf", false, "use fzf interface on results")
		ListCmd.Flags().Bool("fzf-internal", false, "use fzf interface on results")
		ListCmd.Flags().MarkHidden("fzf-internal")
	}
	ListCmd.RegisterFlagCompletionFunc("target", util.BranchCompletion)
}
