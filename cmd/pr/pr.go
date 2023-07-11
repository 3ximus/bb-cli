package pr

import (
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/remote"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

var PrCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlag("repo", cmd.Flags().Lookup("repo"))
		cobra.CheckErr(err)
		if curRepo := getCurrRepo(); curRepo != "" {
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

func getCurrRepo() string {
	url, err := git.Remote(remote.GetURL("origin"))
	if err != nil {
		return ""
	}
	// remotePattern, err := regexp.Compile(`git@github.com:([^\.]*/[^\.]*).git`)
	remotePattern, err := regexp.Compile(`git@bitbucket.org:([^\.]*/[^\.]*).git`)
	if err != nil {
		return ""
	}
	if !remotePattern.MatchString(url) {
		return ""
	}
	return remotePattern.ReplaceAllString(strings.Trim(url, "\n"), "$1")
}
