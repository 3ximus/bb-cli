package cmd

import (
	"bb/cmd/auth"
	"bb/cmd/pr"
	"os"
	"regexp"
	"strings"

	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/remote"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "bb",
	Short: "bb is a bitbucket cli",
	Long:  `Bitbucket cli to interact with bitbucket.org`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// globally set config path
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/bb.yaml)")
	// globally set the repository to use
	rootCmd.PersistentFlags().StringP("repository", "r", "", "selected repository")

	viper.SetDefault("api", "https://api.bitbucket.org/2.0")

	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(pr.PrCmd)
}

func getCurrRepo() string {
	url, err := git.Remote(remote.GetURL("origin"))
	if err != nil {
		return ""
	}
	// remotePattern, err := regexp.Compile(`git@github.com:([^\.]*/[^\.]*).git`)
	remotePattern, err := regexp.Compile(`git@github.com:([^\.]*/[^\.]*).git`)
	if err != nil {
		return ""
	}
	if !remotePattern.MatchString(url) {
		return ""
	}
	return remotePattern.ReplaceAllString(strings.Trim(url, "\n"), "$1")
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

	// repository setup
	err = viper.BindPFlag("repository", rootCmd.Flags().Lookup("repository"))
	if curRepo := getCurrRepo(); curRepo != "" {
		viper.SetDefault("repository", curRepo)
	}
	if !viper.IsSet("repository") {
		cobra.CheckErr("repository is not defined")
	}

	// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
}
