package auth

import (
	"bb/api"
	"bb/util"
	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of your authentication settings.",
	Run: func(cmd *cobra.Command, args []string) {

		user := api.GetUser()

		util.Table([]string{"ID", "Username", "Name", "Link"}, [][]string{
			{user.AccountId, user.Username, user.DisplayName, user.Links.Html.Href},
		})
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
