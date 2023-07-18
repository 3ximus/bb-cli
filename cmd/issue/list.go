package issue

import (
	"bb/api"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	Run: func(cmd *cobra.Command, args []string) {

		cc := api.GetIssueList(viper.GetString("repo"))
		<-cc

		fmt.Println("Not implemented")

	},
}

func init() {
}
