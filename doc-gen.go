package doc

import (
	"bb/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	doc.GenManTree(cmd.RootCmd, nil, "doc")
}
