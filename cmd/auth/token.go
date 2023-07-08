package auth

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Outputs your bitbucket token",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.GetString("token"))
	},
}

func init() {
}
