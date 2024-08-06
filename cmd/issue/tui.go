package issue

import (
	"bb/api"
	"bb/util"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"

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
type item api.JiraIssue

func (i item) FilterValue() string { return i.Key + i.Fields.Summary }

// List item delegate defines how to render a list
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, _ := listItem.(item) // convert to api.JiraIssue
	fmtStr := "  %s \033[1;32m%s\033[m %s %s %s"
	if index == m.Index() { // selected index
		fmtStr = "\033[1;35m| %s \033[1;35m%s %s \033[1;35m%s %s"
	}
	fmt.Fprintf(w, fmtStr,
		util.FormatIssueStatus(i.Fields.Status.Name),
		i.Key,
		util.FormatIssueType(i.Fields.Type.Name),
		i.Fields.Summary,
		util.FormatIssuePriority(i.Fields.Priority.Id))
}

// ====================================================
//                  KEYBINDINGS
// ====================================================

type listKeyMap struct {
	quit          key.Binding
	toggleSpinner key.Binding
	openWeb       key.Binding
	transition    key.Binding
	refresh       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		openWeb: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "open in browser"),
		),
		transition: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "transition issue"),
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
type issueList struct {
	data []list.Item
}

func getIssueList() tea.Msg {
	items := []list.Item{}
	for issue := range api.GetIssueList(100, false, false, "", []string{}, []string{}, "", false) {
		items = append(items, item(issue))
	}
	return issueList{items}
}

// Function that will just trigger the startup update
func startup() tea.Msg {
	return triggerSpinner(true)
}

// ====================================================
//                     COMMANDS
// ====================================================

func loadIssues() tea.Cmd {
	return tea.Batch(startup, getIssueList)
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
	return loadIssues()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(int(math.Min(float64(msg.Height), 20)))
		return m, nil
	case triggerSpinner:
		return m, m.list.StartSpinner()
	case issueList:
		return m, tea.Batch(
			m.list.SetItems(msg.data),
			m.list.ToggleSpinner(),
		)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.toggleSpinner):
			return m, m.list.ToggleSpinner()
		case key.Matches(msg, m.keys.refresh):
			return m, loadIssues()
		case key.Matches(msg, m.keys.openWeb):
			i, ok := m.list.SelectedItem().(item)
			if ok {
				util.OpenInBrowser(api.JiraBrowse(viper.GetString("jira_domain"), i.Key))
				return m, nil
			}
		case key.Matches(msg, m.keys.transition):
			i, ok := m.list.SelectedItem().(item)
			if ok {
				ex, err := os.Executable()
				if err != nil {
					panic(err)
				}
				return m, tea.Sequence(
					tea.ExecProcess(exec.Command(ex, "issue", "transition", i.Key), nil),
					loadIssues())
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
	l.Title = "Jira issue list:"
	l.InfiniteScrolling = true
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(4))
	l.SetShowStatusBar(false)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(2).PaddingTop(0)

	listKeys := newListKeyMap()
	l.AdditionalFullHelpKeys =
		func() []key.Binding {
			return []key.Binding{
				listKeys.quit,
				listKeys.toggleSpinner,
				listKeys.transition,
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
