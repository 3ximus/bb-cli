package pr

import (
	"bb/api"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Status string

const (
	OPEN       Status = "open"
	MERGED     Status = "merged"
	DECLINED   Status = "declined"
	SUPERSEDED Status = "superseded"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *Status) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *Status) Set(v string) error {
	switch v {
	case "open", "merged", "declined", "superseded":
		*e = Status(v)
		return nil
	default:
		return errors.New(`must be one of "open", "merged", "declined" or "superseded"`)
	}
}

// Type is only used in help text
func (e *Status) Type() string {
	return "Status"
}

var state = OPEN

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests from a repository",
	Run: func(cmd *cobra.Command, args []string) {
		author, _ := cmd.Flags().GetString("author")
		search, _ := cmd.Flags().GetString("search")
		fmt.Println(search)
		prs := api.GetPr(viper.GetString("repo"), strings.ToUpper(state.String()), author, search)

		os.Stdout.WriteString("\n")
		for _, pr := range prs {
			// if we didn't provide filter don't show the pr status
			fmt.Printf("%s \033[1;32m#%d\033[m %s  \033[1;34m[ %s → %s]\033[m\n", formatState(pr.State), pr.ID, pr.Title, pr.Source.Branch.Name, pr.Destination.Branch.Name)
			fmt.Printf("%s\033[33m%s\033[m\n", strings.Repeat(" ", len(formatState(pr.State))-4), pr.Author.Nickname)
		}
		os.Stdout.WriteString("\n")
	},
}

func init() {
	ListCmd.Flags().StringP("author", "a", "", "Filter by author nick name")
	ListCmd.Flags().StringP("search", "S", "", "Search pull request with query")
	ListCmd.Flags().VarP(&state, "state", "s", `Filter by state. Default: "open"
Possible options: "open", "merged", "declined" or "superseded"`)
}

func formatState(state string) string {
	stateString := ""
	switch state {
	case "OPEN":
		stateString = "\033[1;44m OPEN \033[m"
	case "MERGED":
		stateString = "\033[1;45m MERGED \033[m"
	case "DECLINED":
		stateString = "\033[1;41m DECLINED \033[m"
	case "SUPERSEDED":
		stateString = "\033[1;44m SUPERSEDED \033[m"
	}
	return stateString
}
