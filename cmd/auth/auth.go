package auth

import (
	"fmt"
	"github.com/spf13/cobra"
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication settings and configuration",
	Long: `This command can be used to setup your authentication with bitbucket.`,
}

func init() {
	AuthCmd.AddCommand(tokenCmd)
	AuthCmd.AddCommand(StatusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
