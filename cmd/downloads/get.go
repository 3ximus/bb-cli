package downloads

import (
	"bb/api"
	"bb/util"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GetCmd = &cobra.Command{
	Use:   "get [ID]",
	Short: "Get a file from the repository downloads",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")
		outputFile, _ := cmd.Flags().GetString("output")
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
		util.Printf("Downloading %s...\n", fileToGet)

		path, err := api.GetDownloadItem(repo, fileToGet, outputFile)
		cobra.CheckErr(err)
		util.Printf("File downloaded: %s\n", path)
	},
}

func init() {
	GetCmd.Flags().StringP("output", "o", "", "Output file name")
	GetCmd.Flags().Bool("latest", false, "Get the latest file")
}
