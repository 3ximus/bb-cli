package pr

import (
	"bb/api"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests from a repository",
	Run: func(cmd *cobra.Command, args []string) {
		stateFilter := genStateFilter(cmd)
		prs := api.GetPr(viper.GetString("repo"), stateFilter, []string{})

		os.Stdout.WriteString("\n")
		for _, pr := range prs {
			// if we didn't provide filter don't show the pr status
			state := ""
			if openFlag, _ := cmd.Flags().GetBool("open"); (len(stateFilter) == 1 && stateFilter[0] != "OPEN") || openFlag {
				state = formatState(pr.State) + "  "
			}
			fmt.Printf(" \033[1;32m#%d\033[m %s%s  \033[1;34m[ %s → %s]\033[m\n", pr.ID, state, pr.Title, pr.Source.Branch.Name, pr.Destination.Branch.Name)
			fmt.Printf("      \033[33m%s\033[m\n", pr.Author.DisplayName)
		}
		os.Stdout.WriteString("\n")
	},
}

func init() {
	ListCmd.Flags().Bool("open", false, "Filter open pull requests")
	ListCmd.Flags().Bool("merged", false, "Filter merged pull requests")
	ListCmd.Flags().Bool("declined", false, "Filter declined pull requests")
	ListCmd.Flags().Bool("superseded", false, "Filter superseded pull requests")

	// ListCmd.Flags().Strin
}

func genStateFilter(cmd *cobra.Command) []string {
	stateFilter := []string{}
	if openFlag, _ := cmd.Flags().GetBool("open"); openFlag {
		stateFilter = append(stateFilter, "OPEN")
	}
	if mergedFlag, _ := cmd.Flags().GetBool("merged"); mergedFlag {
		stateFilter = append(stateFilter, "MERGED")
	}
	if declinedFlag, _ := cmd.Flags().GetBool("declined"); declinedFlag {
		stateFilter = append(stateFilter, "DECLINED")
	}
	if supersededFlag, _ := cmd.Flags().GetBool("superseded"); supersededFlag {
		stateFilter = append(stateFilter, "SUPERSEDED")
	}
	if len(stateFilter) == 0 {
		stateFilter = append(stateFilter, "OPEN")
	}
	return stateFilter
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
