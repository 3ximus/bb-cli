package downloads

import (
	"bb/api"
	"bb/util"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete [ID]",
	Short: "Delete a file from the repository downloads",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")
		getLatest, _ := cmd.Flags().GetBool("latest")

		var fileToGet string
		if getLatest {
			fileToGet = (<-api.GetDownloadsList(repo)).Name
		} else {
			if len(args) != 1 {
				fmt.Errorf("Needs an argument or --latest flag")
				return
			}
			fileToGet = args[0]
		}

		util.Printf("Delete file %s ? [y/n]\n", fileToGet)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if strings.TrimSpace(strings.ToLower(scanner.Text())) != "y" {
			return
		}

		api.DeleteDownloadItem(repo, fileToGet)
		util.Printf("File deleted")
	},
}

func init() {
	DeleteCmd.Flags().Bool("latest", false, "Get the latest file")
}
