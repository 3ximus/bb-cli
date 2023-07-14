# Bitbucket CLI

A cli tool for bitbucket similar to [gh](https://cli.github.com/) written in `go` for bitbucket.org API

## Instalation

```bash
go install -ldflags="-s -w"
```

### Setup autocompletion

For example for bash

```bash
bb completion bash > bb.bash
sudo mv bb.bash /usr/share/bash-completion/completions/
```

## Usage

TODO





### NOTE TO SELF

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
