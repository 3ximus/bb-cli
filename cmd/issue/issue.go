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

		err = viper.BindPFlag("jira_api", cmd.Flags().Lookup("endpoint"))
		cobra.CheckErr(err)
		if !viper.IsSet("jira_api") {
			cobra.CheckErr("jira endpoint is not defined")
		}
	},
}

func init() {
	IssueCmd.AddCommand(ListCmd)
	IssueCmd.AddCommand(ViewCmd)
	IssueCmd.PersistentFlags().StringP("repo", "R", "", "selected repository")
	IssueCmd.PersistentFlags().StringP("endpoint", "E", "", `endpoint for your organization api on jira.
	Format: https://XXXXXX.atlassian.net/rest/api/3
	`)
}
