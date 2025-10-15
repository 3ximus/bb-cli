
package tempo

import (
	"bb/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var TempoCmd = &cobra.Command{
	Use:   "tempo",
	Short: "Manage tempo worklogs",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlag("repo", cmd.Flags().Lookup("repo"))
		cobra.CheckErr(err)
		if curRepo := util.GetCurrentRepo(); curRepo != "" {
			viper.SetDefault("repo", curRepo)
		}

		err = viper.BindPFlag("jira_domain", cmd.Flags().Lookup("domain"))
		cobra.CheckErr(err)
		if !viper.IsSet("jira_domain") {
			cobra.CheckErr("jira domain is not defined")
		}
	},
}

func init() {
	TempoCmd.AddCommand(ListCmd)
	TempoCmd.PersistentFlags().StringP("repo", "R", "", "selected repository")
	TempoCmd.PersistentFlags().StringP("domain", "D", "", "your jira domain ( XXXX in https://XXXX.atlassian.net)")
}
