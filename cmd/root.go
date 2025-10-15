package cmd

import (
	"bb/cmd/auth"
	"bb/cmd/doc"
	"bb/cmd/downloads"
	"bb/cmd/environment"
	"bb/cmd/issue"
	"bb/cmd/pipeline"
	"bb/cmd/pr"
	"bb/cmd/tempo"
	"bb/store"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "bb",
	Short: "CLI utility to manage Bitbucket repositories and Jira organizations",
	Long: `This utility is focused on allowing simple operations on bitbucket and jira through the command line.
	It provides commands to operate on Bitbucket and Jira.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// globally set config path
	RootCmd.PersistentFlags().StringVar(&store.CfgFile, "config", "", "config file (default is $HOME/.config/bb.yaml)")
	RootCmd.PersistentFlags().BoolVar(&store.UseColor, "color", false, "use color even if stdout is piped")

	RootCmd.AddCommand(auth.AuthCmd)
	RootCmd.AddCommand(pr.PrCmd)
	RootCmd.AddCommand(environment.EnvironmentCmd)
	RootCmd.AddCommand(issue.IssueCmd)
	RootCmd.AddCommand(tempo.TempoCmd)
	RootCmd.AddCommand(pipeline.PipelineCmd)
	RootCmd.AddCommand(downloads.DownloadsCmd)
	RootCmd.AddCommand(doc.DocCmd)
}

func initConfig() {
	if store.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(store.CfgFile)
	} else {
		configDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		// Search config in current directory or in .config
		viper.AddConfigPath(configDir)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("bb")
	}
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	cobra.CheckErr(err)

	viper.SetDefault("bb_api", "https://api.bitbucket.org/2.0")
	viper.SetDefault("tempo_api", "https://api.tempo.io/4")
}
