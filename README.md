# Attlassian CLI

A cli tool for bitbucket and jira similar to [gh](https://cli.github.com/) written in `go`

_Currently under development_

## [Documentation](https://github.com/3ximus/bb-cli/wiki/bb)

The full documentation is on the wiki of this project:

https://github.com/3ximus/bb-cli/wiki/bb

## Instalation

```bash
go install -ldflags="-s -w"
```

## Configuration

This is an example of possible config options:
A full example can be seen in `bb.example.yaml`

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

# define custom text icons or text for jira_status, jira_types, pr_status or pipeline_status. The format is as follows:
#   identifier:
#     values: ["State 1", "State 2"] # this is the string that matches the state being printed
#     color: "1;34" # the ANSII sequence for the color used. if I want to remove icon 1;38;5;235;44
#     text: "ACTUAL STATE" # string printed, takes precedence over icon
#     icon: "ﲊ" # icon to display
#
# examples bellow:

jira_status:
  inprogress:
    values: ["In Progress", "In Progress_T"]
    icon: ""
    color: "1;34"
  todo:
    values: ["À FAIRE"]
    icon: ""
    color: "1;33"
```

### Setup autocompletion

Generate completion for your shell with `bb completion <your-shell>` and save the content in your completions directory

For example for bash

```bash
bb completion bash | sudo tee /usr/share/bash-completion/completions/bb.bash >/dev/null
```

## Usage

Use the help command to display usage information for each command

```
bb help
bb help [COMMAND]
```

## TODO
