package pr

import (
	"github.com/spf13/cobra"
)

var PrCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	PrCmd.AddCommand(ListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
