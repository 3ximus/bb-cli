package cmd

import (
	"bb/cmd/auth"
	"bb/cmd/pr"
	"bb/cmd/issue"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "bb",
	Short: "bb is a bitbucket cli",
	Long:  `Bitbucket cli to interact with bitbucket.org`,
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
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/bb.yaml)")

	viper.SetDefault("bb_api", "https://api.bitbucket.org/2.0")

	RootCmd.AddCommand(auth.AuthCmd)
	RootCmd.AddCommand(pr.PrCmd)
	RootCmd.AddCommand(issue.IssueCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
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

	// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
}
