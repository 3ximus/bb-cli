package issue

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented")

	},
}

func init() {
}
