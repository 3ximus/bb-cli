# Bitbucket CLI

A cli tool for bitbucket similar to [gh](https://cli.github.com/) written in `go` for bitbucket.org API




## NOTE

To generate documentation use this

```go
package doc

import (
	"bb/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	doc.GenManTree(cmd.RootCmd, nil, "doc")
}
```
