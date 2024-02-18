package pipeline

import (
	"bb/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"pl"},
	Short:   "Manage pipelines",
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
	PipelineCmd.AddCommand(ListCmd)
	PipelineCmd.AddCommand(ViewCmd)
	PipelineCmd.AddCommand(StopCmd)
	PipelineCmd.AddCommand(RunCmd)
	PipelineCmd.PersistentFlags().StringP("repo", "R", "", "selected repository")
}
