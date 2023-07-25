package util

import (
	"bb/api"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func FormatPrState(state api.PrState) string {
	stateString := ""
	switch state {
	case "OPEN":
		stateString = "\033[1;38;5;235;44m OPEN \033[m"
	case "MERGED":
		stateString = "\033[1;38;5;235;45m MERGED \033[m"
	case "DECLINED":
		stateString = "\033[1;38;5;235;41m DECLINED \033[m"
	case "SUPERSEDED":
		stateString = "\033[1;38;5;235;44m SUPERSEDED \033[m"
	}
	return stateString
}

func FormatPipelineState(state string) string {
	stateString := ""
	switch state {
	case "INPROGRESS":
		stateString = "\033[1;38;5;235;44m RUNNING \033[m"
	case "STOPPED", "merged":
		stateString = "\033[1;38;5;235;43m STOPPED \033[m"
	case "SUCCESSFUL", "declined":
		stateString = "\033[1;38;5;235;42m SUCCESSFUL \033[m"
	case "FAILED", "superseded":
		stateString = "\033[1;38;5;235;41m FAILED \033[m"
	}
	return stateString
}

func FormatIssueStatus(status string) string {
	color := "\033[1;38;5;235;47m"
	statusMap := viper.GetStringMapStringSlice("jira_status")
	for k, v := range statusMap {
		for _, s := range v {
			if s == status {
				switch strings.ToUpper(k) {
				// case "OPEN":
				// 	color = "\033[1;38;5;235;44m"
				case "INPROGRESS":
					color = "\033[1;38;5;235;44m"
				case "TODO":
					color = "\033[1;38;5;235;43m"
				case "TESTING":
					color = "\033[1;38;5;235;46m"
				case "DONE":
					color = "\033[1;38;5;235;42m"
				case "BLOCKED":
					color = "\033[1;38;5;235;41m"
				}
			}
		}
	}
	return fmt.Sprintf("%s %s \033[m", color, status)
}

func FormatIssuePriority(id string, name string) string {
	priorityString := ""
	switch id {
	case "1":
		priorityString = fmt.Sprintf("\033[1;31m▲ %s\033[m", name)
	case "2":
		priorityString = fmt.Sprintf("\033[1;35m▲ %s\033[m", name)
	case "3":
		priorityString = fmt.Sprintf("\033[1;33m⯀ %s\033[m", name)
	case "4":
		priorityString = fmt.Sprintf("\033[1;34m▼ %s\033[m", name)
	default:
		priorityString = fmt.Sprintf("\033[1;37m%s\033[m", name)
	}
	return priorityString
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

func OpenInBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	cobra.CheckErr(err)
}

func UseExternalFZF[T any](list []T, prompt string, toString func(int) string) []int {
	input := ""
	for i := range list {
		input += fmt.Sprintf("%d %s\n", i, toString(i))
	}
	cmd := exec.Command("fzf", "-m", "--height", "20%", "--reverse", "--with-nth", "2..", "--prompt", prompt)
	var selectionBuffer strings.Builder
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &selectionBuffer
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	cobra.CheckErr(err)
	err = cmd.Wait()

	var result []int
	for _, r := range strings.Split(selectionBuffer.String(), "\n") {
		if r == "" {
			continue
		}
		idx, err := strconv.Atoi(strings.Split(r, " ")[0])
		cobra.CheckErr(err)
		result = append(result, idx)
	}
	return result
}
