package downloads

import (
	"bb/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DownloadsCmd = &cobra.Command{
	Use:     "downloads",
	Aliases: []string{"dl"},
	Short:   "Manage downloads [dl]",
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
	DownloadsCmd.AddCommand(ListCmd)
	DownloadsCmd.AddCommand(GetCmd)
	DownloadsCmd.AddCommand(DeleteCmd)
	DownloadsCmd.PersistentFlags().StringP("repo", "R", "", "selected repository")
}
