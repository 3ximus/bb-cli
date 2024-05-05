package issue

import (
	"bb/api"
	"bb/util"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		util.FormatIssuePriority(i.Fields.Priority.Id, i.Fields.Priority.Name))
}

// ====================================================
//                     MESSAGES
// ====================================================

type started bool
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
// turning on the spinner
func startup() tea.Msg {
	return started(true)
}

// ====================================================
//                   TUI METHODS
// ====================================================

// The main TUI model
type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return tea.Batch(startup, getIssueList)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case started:
		return m, m.list.StartSpinner()
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case issueList:
		return m, tea.Batch(
			m.list.SetItems(msg.data),
			m.list.ToggleSpinner(),
		)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "s":
			return m, m.list.ToggleSpinner()
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				return m, tea.Batch(
					tea.Printf("ENTER ON %s", i.Key),
				)
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
	l := list.New([]list.Item{}, itemDelegate{}, 0, 13)
	l.Title = "Jira issue list:"
	l.InfiniteScrolling = true
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(4))
	l.SetShowStatusBar(false)
	// TODO
	// l.Styles.Spinner = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(1))
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(2).PaddingTop(0)

	m := model{list: l}

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
