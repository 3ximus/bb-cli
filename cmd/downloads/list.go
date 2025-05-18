package downloads

import (
	"bb/api"
	"bb/util"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List downloads from a repository",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		useFZF, _ := cmd.Flags().GetBool("fzf")
		useFZFInternal, _ := cmd.Flags().GetBool("fzf-internal")
		if useFZF {
			util.ReplaceListWithFzf("--read0 --prompt 'View > '" +
				" --header='\033[1;33mctrl-g\033[m: get file | \033[1;33mctrl-d\033[m: delete file'" +
				" --bind 'enter:become(" + os.Args[0] + " downloads -R " + viper.GetString("repo") + " get {2})'" +
				" --bind 'ctrl-g:execute(" + os.Args[0] + " downloads -R " + viper.GetString("repo") + " get {2})'" +
				" --bind 'ctrl-e:execute(" + os.Args[0] + " downloads -R " + viper.GetString("repo") + " delete {2}")
			return
		}

		count := 0
		for downloadItem := range api.GetDownloadsList(viper.GetString("repo")) {
			util.Printf("\033[1;33m%s\033[m  %s  \033[37m(%s)\033[m", util.FormatBytes(downloadItem.Size), downloadItem.Name, util.TimeAgo(downloadItem.CreatedOn))

			endChar := "\n"
			if useFZFInternal {
				endChar = "\x00"
			}
			util.Printf(endChar)

			count++
		}
		if count == 0 {
			util.Printf("No downloads for \033[1;36m%s\033[m\n", viper.GetString("repo"))
		}

	},
}

func init() {
	// ListCmd.Flags().IntP("number-results", "n", 10, "max number of results retrieve (max: 100)")
	if util.CommandExists("fzf") {
		ListCmd.Flags().Bool("fzf", false, "use fzf interface on results")
		ListCmd.Flags().Bool("fzf-internal", false, "use fzf interface on results")
		ListCmd.Flags().MarkHidden("fzf-internal")
	}
}
