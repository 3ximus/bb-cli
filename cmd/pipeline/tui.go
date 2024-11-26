package pipeline

import (
	"bb/api"
	"bb/util"
	"fmt"
	"io"
	"math"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

// ====================================================
//                     ITEMS
// ====================================================

// Defines the structure of each item in a list
type item api.Pipeline

func (i item) FilterValue() string { return i.State.Name }

// List item delegate defines how to render a list
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, _ := listItem.(item) // convert to api.Pipeline
	fmtStr := "  %s \033[1;32m#%d\033[m %s \033[37m%s (%s)\033[m"
	if index == m.Index() { // selected index
		fmtStr = "\033[1;35m|  %s \033[1;35m#%d %s \033[1;35m%s (%s)\033[m"
	}

	status := ""
	if i.State.Result.Name == "" {
		status = fmt.Sprintf("%s", util.FormatPipelineStatus(i.State.Name))
	} else {
		status = fmt.Sprintf("%s", util.FormatPipelineStatus(i.State.Result.Name))
	}

	source := ""
	if i.Target.Source != "" {
		source = fmt.Sprintf("%s \033[1;34m[ %s → %s ]\033[m", i.Target.PullRequest.Title, i.Target.Source, i.Target.Destination)
	} else {
		source = fmt.Sprintf("\033[1;34m[ %s ]\033[m", i.Target.RefName)
	}

	// fmt.Fprintf(w, " \033[37m%s (%s)\033[m", util.TimeDuration(time.Duration(i.DurationInSeconds*1000000000)), util.TimeAgo(i.CreatedOn))

	fmt.Fprintf(w, fmtStr,
		status,
		i.BuildNumber,
		source,
		util.TimeDuration(time.Duration(i.DurationInSeconds*1000000000)), util.TimeAgo(i.CreatedOn))

}

// ====================================================
//                  KEYBINDINGS
// ====================================================

type listKeyMap struct {
	quit    key.Binding
	openWeb key.Binding
	stop    key.Binding
	refresh key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		openWeb: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "open in browser"),
		),
		stop: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "stop pipeline"),
		),
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh list"),
		),
	}
}

// ====================================================
//                     MESSAGES
// ====================================================

type triggerSpinner bool
type updatedList struct {
	data []list.Item
}

func getPipelineList() tea.Msg {
	items := []list.Item{}
	for pipeline := range api.GetPipelineList(viper.GetString("repo"), 10, "") {
		items = append(items, item(pipeline))
	}
	return updatedList{items}
}

// Function that will just trigger the startup update
func startup() tea.Msg {
	return triggerSpinner(true)
}

// ====================================================
//                     COMMANDS
// ====================================================

func loadPipelines() tea.Cmd {
	return tea.Batch(startup, getPipelineList)
}

// ====================================================
//                   TUI METHODS
// ====================================================

// The main TUI model
type model struct {
	list list.Model
	keys *listKeyMap
}

func (m model) Init() tea.Cmd {
	return loadPipelines()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(int(math.Min(float64(msg.Height), 10)))
		return m, nil
	case triggerSpinner:
		return m, m.list.StartSpinner()
	case updatedList:
		return m, tea.Batch(
			m.list.SetItems(msg.data),
			m.list.ToggleSpinner(),
		)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.refresh):
			return m, loadPipelines()
		case key.Matches(msg, m.keys.openWeb):
			i, ok := m.list.SelectedItem().(item)
			if ok {
				util.OpenInBrowser(api.BBBrowsePipelines(viper.GetString("repo"), i.BuildNumber))
				return m, nil
			}
		case key.Matches(msg, m.keys.stop):
			i, ok := m.list.SelectedItem().(item)
			if ok {
				api.StopPipeline(viper.GetString("repo"), fmt.Sprintf("%d", i.BuildNumber))
				return m, loadPipelines()
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}

// Main Function
func TUI() {
	l := list.New([]list.Item{}, itemDelegate{}, 0, 10)
	l.Title = "Pipeline list:"
	l.InfiniteScrolling = true
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(4))
	l.SetShowStatusBar(false)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(2).PaddingTop(0)

	listKeys := newListKeyMap()
	l.AdditionalFullHelpKeys =
		func() []key.Binding {
			return []key.Binding{
				listKeys.quit,
				listKeys.stop,
				listKeys.openWeb,
				listKeys.refresh,
			}
		}

	m := model{list: l, keys: listKeys}

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
