package pr

import (
	"bb/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PrCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests",
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
	PrCmd.AddCommand(ListCmd)
	PrCmd.AddCommand(CreateCmd)
	PrCmd.AddCommand(ViewCmd)
	PrCmd.PersistentFlags().StringP("repo", "R", "", "selected repository")
}
