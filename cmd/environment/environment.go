package environment

import (
	"bb/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var EnvironmentCmd = &cobra.Command{
	Use:     "environment",
	Aliases: []string{"env"},
	Short:   "Manage environments [env]",
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
	EnvironmentCmd.AddCommand(ListCmd)
	EnvironmentCmd.PersistentFlags().StringP("repo", "R", "", "selected repository")
}
