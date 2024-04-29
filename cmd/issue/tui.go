package issue

import (
	"bb/api"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// create global styles
var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("3"))
)

var listingChannel <-chan api.JiraIssue

type newDataLoaded struct {
	data []list.Item
}

// Defines the structure of each item in a list
type item api.JiraIssue

func (i item) Title() string       { return i.Key }
func (i item) Description() string { return i.Fields.Summary }
func (i item) FilterValue() string { return i.Key }

// List item delegate defines how to render a list
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Key)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("| " + strings.Join(s, ""))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list       list.Model
	data       []api.JiraIssue
	dataLoaded bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

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
	case newDataLoaded:
		return m, m.list.SetItems(msg.data)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}

func TUI() {
	listingChannel = api.GetIssueList(100, false, false, "", []string{}, []string{}, "", false)

	l := list.New([]list.Item{}, itemDelegate{}, 0, 13)
	l.Title = "Issue list:"
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(4))
	l.SetShowStatusBar(false)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(2).PaddingTop(0)

	m := model{list: l, dataLoaded: false}

	p := tea.NewProgram(m)

	// load data async into the list
	go func() {
		items := []list.Item{}
		for issue := range listingChannel {
			items = append(items, item{Key: issue.Key})
		}
		p.Send(newDataLoaded{items})
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
