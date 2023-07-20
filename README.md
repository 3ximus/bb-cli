# Bitbucket CLI

A cli tool for bitbucket and jira similar to [gh](https://cli.github.com/) written in `go` for bitbucket.org API

## Instalation

```bash
go install -ldflags="-s -w"
```

## Configuration

This is an example of possible config options:

`$HOME/.config/bb.yaml` or `./bb.yaml`

```yaml
# Authentication options for Bitbucket
bb_token: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
username: fabio_almeida_vo2

# Authentication options for Jira
jira_domain: XXXXXXXXX # In https://<your-domain>.atlassian.net
email: your@email.com
jira_token: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX(192 characters)

# Extra options:

# include branch name at beggining of the pull request (useful to link with jira tickets)
include_branch_name: true
```

### Setup autocompletion

Generate completion for your shell with `bb completion <your-shell>` and save the content in your completions directory

For example for bash

```bash
bb completion bash | sudo tee /usr/share/bash-completion/completions/bb.bash >/dev/null
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
