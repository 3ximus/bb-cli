package pipeline

import (
	"bb/api"
	"bb/util"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running pipeline",
	Args:  cobra.MaximumNArgs(1),
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

		api.StopPipeline(repo, fmt.Sprintf("%d", id))
		fmt.Printf("Pipeline #%d \033[1;31mStopped\033[m\n", id)
	},
}

func init() {
}
