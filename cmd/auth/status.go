package auth

import (
	"bb/api"
	"fmt"
	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of your authentication settings.",
	Run: func(cmd *cobra.Command, args []string) {

		user := api.GetUser()

		fmt.Printf("\n \033[1;34mID\033[m       %s\n \033[1;34mUsername\033[m %s\n \033[1;34mName\033[m     %s\n \033[1;34mLink\033[m     %s\n\n",
			user.AccountId, user.Username, user.DisplayName, user.Links.Html.Href)
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
