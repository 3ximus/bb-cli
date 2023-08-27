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

	"github.com/ktr0731/go-fuzzyfinder"
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
	case "STOPPED", "stopped":
		stateString = "\033[1;38;5;235;43m STOPPED \033[m"
	case "SUCCESSFUL", "successful":
		stateString = "\033[1;38;5;235;42m SUCCESSFUL \033[m"
	case "FAILED", "failed":
		stateString = "\033[1;38;5;235;41m FAILED \033[m"
	}
	return stateString
}

func FormatIssueType(issueType string) string {
	str := "\033[1;37m"
	typeMap := viper.GetStringMapStringSlice("jira_type")
	iconMap := viper.GetStringMap("jira_type_icon")
	for k, v := range typeMap {
		for _, t := range v {
			if t == issueType {
				switch strings.ToUpper(k) {
				case "BUG":
					if icon, ok := iconMap[k]; ok {
						str = fmt.Sprintf("\033[1;31m%s", icon)
					} else {
						str = fmt.Sprintf("\033[1;31m%s", issueType)
					}
				case "TASK":
					if icon, ok := iconMap[k]; ok {
						str = fmt.Sprintf("\033[1;34m%s", icon)
					} else {
						str = fmt.Sprintf("\033[1;34m%s", issueType)
					}
				case "EPIC":
					if icon, ok := iconMap[k]; ok {
						str = fmt.Sprintf("\033[1;35m%s", icon)
					} else {
						str = fmt.Sprintf("\033[1;35m%s", issueType)
					}
				}
			}
		}
	}
	return fmt.Sprintf("%s\033[m", str)
}

func FormatIssueStatus(status string) string {
	str := fmt.Sprintf("\033[1;38;5;235;47m %s \033[m", status) // default
	statusMap := viper.GetStringMapStringSlice("jira_status")
	iconMap := viper.GetStringMap("jira_status_icons")
	for k, v := range statusMap {
		for _, s := range v {
			if s == status {
				switch strings.ToUpper(k) {
				case "OPEN":
					if icon, ok := iconMap[k]; ok {
						str = fmt.Sprintf("\033[1;34m%s", icon)
					} else {
						str = fmt.Sprintf("\033[1;38;5;235;47m %s \033[m", status)
					}
				case "INPROGRESS":
					if icon, ok := iconMap[k]; ok {
						str = fmt.Sprintf("\033[1;34m%s", icon)
					} else {
						str = fmt.Sprintf("\033[1;38;5;235;44m %s \033[m", status)
					}
				case "TODO":
					if icon, ok := iconMap[k]; ok {
						str = fmt.Sprintf("\033[1;33m%s", icon)
					} else {
						str = fmt.Sprintf("\033[1;38;5;235;43m %s \033[m", status)
					}
				case "TESTING":
					if icon, ok := iconMap[k]; ok {
						str = fmt.Sprintf("\033[1;36m%s", icon)
					} else {
						str = fmt.Sprintf("\033[1;38;5;235;46m %s \033[m", status)
					}
				case "DONE":
					if icon, ok := iconMap[k]; ok {
						str = fmt.Sprintf("\033[1;32m%s", icon)
					} else {
						str = fmt.Sprintf("\033[1;38;5;235;42m %s \033[m", status)
					}
				case "BLOCKED":
					if icon, ok := iconMap[k]; ok {
						str = fmt.Sprintf("\033[1;31m%s", icon)
					} else {
						str = fmt.Sprintf("\033[1;38;5;235;41m %s \033[m", status)
					}
				}
			}
		}
	}
	return str
}

func FormatIssuePriority(id string, name string) string {
	priorityString := ""
	switch id {
	case "1":
		priorityString = fmt.Sprintf("\033[1;31m %s\033[m", name)
	case "2":
		priorityString = fmt.Sprintf("\033[1;35m %s\033[m", name)
	case "3":
		priorityString = fmt.Sprintf("\033[1;33m %s\033[m", name)
	case "4":
		priorityString = fmt.Sprintf("\033[1;34m %s\033[m", name)
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

func ConvertToSeconds(timeStrings []string) (int, error) {
	var seconds int
	for _, timeString := range timeStrings {
		timeParts := strings.Split(timeString, " ")
		for _, part := range timeParts {
			if len(part) == 0 {
				continue
			}
			value, unit := part[:len(part)-1], part[len(part)-1]
			switch unit {
			case 'h':
				hours, err := strconv.Atoi(value)
				if err != nil {
					return 0, fmt.Errorf("Invalid format: %s", timeString)
				}
				seconds += hours * 60 * 60
			case 'm':
				minutes, err := strconv.Atoi(value)
				if err != nil {
					return 0, fmt.Errorf("Invalid format: %s", timeString)
				}
				seconds += minutes * 60
			case 'd':
				days, err := strconv.Atoi(value)
				if err != nil {
					return 0, fmt.Errorf("Invalid format: %s", timeString)
				}
				seconds += days * 8 * 60 * 60
			default:
				return 0, fmt.Errorf("Invalid format: %s", timeString)
			}
		}
	}

	return seconds, nil
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

func OpenInEditor(file *os.File) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	cobra.CheckErr(err)
	err = cmd.Wait()
	cobra.CheckErr(err)
}

func SelectFZF[T any](list []T, prompt string, toString func(int) string) []int {
	if len(list) == 0 {
		return []int{}
	}

	var indexes []int

	// check if fzf is installed
	_, existsErr := exec.LookPath("fzf")
	if existsErr == nil {
		indexes = UseExternalFZF(list, prompt, toString)
	} else {
		var err error
		// backup in case fzf is not installed in the system
		indexes, err = fuzzyfinder.FindMulti(list, toString, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop), fuzzyfinder.WithPromptString(prompt))
		cobra.CheckErr(err)
	}
	return indexes
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
