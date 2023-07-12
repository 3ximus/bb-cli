package util

import (
	"regexp"
	"strings"

	"github.com/ldez/go-git-cmd-wrapper/v2/branch"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/remote"
	"github.com/spf13/cobra"
)

func GetCurrentRepo() string {
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

func GetCurrentBranch() string {
	branch, err := git.Branch(branch.ShowCurrent)
	cobra.CheckErr(err)
	return strings.Trim(branch, "\n")
}
