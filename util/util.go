package util

import (
	"bb/api"
	"fmt"
	"time"
)

func FormatPrState(state api.PrState) string {
	stateString := ""
	switch state {
	case "OPEN":
		stateString = "\033[1;44m OPEN \033[m"
	case "MERGED":
		stateString = "\033[1;45m MERGED \033[m"
	case "DECLINED":
		stateString = "\033[1;41m DECLINED \033[m"
	case "SUPERSEDED":
		stateString = "\033[1;44m SUPERSEDED \033[m"
	}
	return stateString
}

func FormatPipelineState(state string) string {
	stateString := ""
	switch state {
	case "INPROGRESS":
		stateString = "\033[1;44m RUNNING \033[m"
	case "STOPPED", "merged":
		stateString = "\033[1;38;5;235;43m STOPPED \033[m"
	case "SUCCESSFUL", "declined":
		stateString = "\033[1;38;5;235;42m SUCCESSFUL \033[m"
	case "FAILED", "superseded":
		stateString = "\033[1;41m FAILED \033[m"
	}
	return stateString
}

func TimeAgo(updatedOn time.Time) string {
	duration := time.Since(updatedOn)
	if duration.Hours() < 1 {
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	} else if duration.Hours() < 24 {
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	} else if duration.Hours() < 48 {
		return "yesterday"
	} else if duration.Hours() < 720 {
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	} else if duration.Hours() < 8760 {
		return fmt.Sprintf("%d months ago", int(duration.Hours()/720))
	}
	return fmt.Sprintf("%d years ago", int(duration.Hours()/8760))
}

