// vim: foldmethod=indent foldnestmax=1

package util

import (
	"bb/api"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

type ResultSwitchConfig struct {
	Values []string `mapstructure:"values"`
	Text   string   `mapstructure:"text"`
	Icon   string   `mapstructure:"icon"`
	Color  string   `mapstructure:"color"`
}

func FormatSwitchConfig(result string, mapping map[string]ResultSwitchConfig) string {
	str := fmt.Sprintf("\033[1;38;5;235;47m %s \033[m", result) // default

	for k, v := range mapping {
		for _, s := range v.Values {
			if s == result {
				if v.Text != "" {
					str = fmt.Sprintf("\033[%sm %s \033[m", v.Color, v.Text)
				} else if v.Icon != "" {
					str = fmt.Sprintf("\033[%sm%s\033[m", v.Color, v.Icon)
				} else {
					str = fmt.Sprintf("\033[%sm %s \033[m", v.Color, strings.ToUpper(k))
				}
			}
		}
	}
	return str
}

/* Returns the config switch config  */
func GetConfig(configKey string, key string) (ResultSwitchConfig, error) {
	mapping := make(map[string]ResultSwitchConfig)
	if err := viper.UnmarshalKey(configKey, &mapping); err != nil {
		cobra.CheckErr(err)
	}
	for k, v := range mapping {
		if k == key {
			return v, nil
		}
	}
	return ResultSwitchConfig{}, errors.New(fmt.Sprintf("Config key not found: '%s'", configKey))
}

func FormatPrState(state api.PrState) string {
	prStatusMap := make(map[string]ResultSwitchConfig)
	if err := viper.UnmarshalKey("pr_status", &prStatusMap); err != nil {
		cobra.CheckErr(err)
	}
	return FormatSwitchConfig(state.String(), prStatusMap)
}

func FormatPipelineStatus(state string) string {
	pipelineStatusMap := make(map[string]ResultSwitchConfig)
	if err := viper.UnmarshalKey("pipeline_status", &pipelineStatusMap); err != nil {
		cobra.CheckErr(err)
	}
	return FormatSwitchConfig(state, pipelineStatusMap)
}

func FormatIssueType(issueType string) string {
	jiraStatusMap := make(map[string]ResultSwitchConfig)
	if err := viper.UnmarshalKey("jira_type", &jiraStatusMap); err != nil {
		cobra.CheckErr(err)
	}
	return FormatSwitchConfig(issueType, jiraStatusMap)
}

func FormatIssueStatus(status string) string {
	jiraStatusMap := make(map[string]ResultSwitchConfig)
	if err := viper.UnmarshalKey("jira_status", &jiraStatusMap); err != nil {
		cobra.CheckErr(err)
	}
	return FormatSwitchConfig(status, jiraStatusMap)
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
	case "5":
		priorityString = fmt.Sprintf("\033[1;38;5;8m %s\033[m", name)
	default:
		priorityString = fmt.Sprintf("\033[1;m%s\033[m", name)
	}
	return priorityString
}

func TimeAgo(updatedOn time.Time) string {
	duration := time.Since(updatedOn)
	return fmt.Sprintf("%s ago", TimeDuration(duration))
}

func TimeDuration(duration time.Duration) string {
	if duration.Hours() < 1 {
		return fmt.Sprintf("%d minutes", int(duration.Minutes()))
	} else if duration.Hours() < 24 {
		return fmt.Sprintf("%d hours", int(duration.Hours()))
	} else if duration.Hours() < 48 {
		return "yesterday"
	} else if duration.Hours() < 720 {
		return fmt.Sprintf("%d days", int(duration.Hours()/24))
	} else if duration.Hours() < 8760 {
		return fmt.Sprintf("%d months", int(duration.Hours()/720))
	}
	return fmt.Sprintf("%d years", int(duration.Hours()/8760))
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

// EXTERNAL OPTIONS

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

// LOG FUNCTIONS

/* fmt.Printf wrapper to remove ANSI colors if stdout is not a terminal */
func Printf(format string, a ...any) {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Printf(format, a...)
	} else {
		ansiColorRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
		fmt.Printf(ansiColorRegex.ReplaceAllString(format, ""), a...)
	}
}

// COMPLETION DIRECTIVES

func BranchCompletion(comd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return ListBranches(), cobra.ShellCompDirectiveDefault
}
