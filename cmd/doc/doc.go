package doc

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var DocCmd = &cobra.Command{
	Use:    "doc",
	Short:  "Generate documentation",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// doc.GenManTree(cmd.Parent(), nil, "/tmp")

		// TODO need to replace all links in see also to remove .md extension
		//   sed -i 's/\.md//' *.md
		doc.GenMarkdownTree(cmd.Parent(), "../bb-cli.wiki")
	},
}

func init() {
}
