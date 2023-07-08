package auth

import (
	"bb/api"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of your authentication settings.",
	Run: func(cmd *cobra.Command, args []string) {

		user := api.GetUser()

		headerFmt := color.New(color.FgHiBlue, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgHiGreen).SprintfFunc()
		tbl := table.New("ID", "Name", "Username")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		tbl.WithPadding(3)
		tbl.AddRow(user.AccountId, user.DisplayName, user.Username)
		tbl.Print()

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
