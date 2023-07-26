package util

import (
	"errors"
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

func GetCurrentBranch() (string, error) {
	output, err := git.Branch(branch.ShowCurrent)
	if err != nil {
		err = errors.New(output) // error message is actually stored
	}
	return strings.Trim(output, "\n"), err
}

func ListBranches() []string {
	branch, err := git.Branch()
	cobra.CheckErr(err)
	var branches = []string{}
	for _, line := range strings.Split(branch, "\n") {
		branches = append(branches, strings.Trim(line, " *"))
	}
	return branches
}
