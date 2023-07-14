package issue

import (
	"bb/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var IssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues / jira tickets",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlag("repo", cmd.Flags().Lookup("repo"))
		cobra.CheckErr(err)
		if curRepo := util.GetCurrentRepo(); curRepo != "" {
			viper.SetDefault("repo", curRepo)
		}
		if !viper.IsSet("repo") {
			cobra.CheckErr("repo is not defined")
		}
	},
}

func init() {
	IssueCmd.AddCommand(ListCmd)
	IssueCmd.PersistentFlags().StringP("repo", "R", "", "selected repository")
}
