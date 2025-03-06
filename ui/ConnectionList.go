package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConnectionList struct {
	focus bool
	table table.Model
}

func NewConnectionList() ConnectionList {
	var columns = []table.Column{
		{Title: "id", Width: 4},
		{Title: "status", Width: 10},
	}

	styles := table.DefaultStyles()
	styles.Selected = styles.Selected.
		Background(lipgloss.Color("240")).
		UnsetForeground()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true)

	return ConnectionList{
		focus: false,
		table: table.New(
			table.WithColumns(columns),
			table.WithRows([]table.Row{
				{"#46436", "closed"},
				{"#46436", "closed"},
				{"#46436", "closed"},
				{"#46436", "closed"},
				{"#46436", "closed"},
				{"#46436", "closed"},
			}),
			table.WithFocused(false),
			table.WithStyles(styles),
		),
	}
}

func (m *ConnectionList) GetFocus() bool {
	return m.focus
}

func (m *ConnectionList) Focus() {
	m.focus = true
	m.table.Focus()
}

func (m *ConnectionList) Blur() {
	m.focus = false
	m.table.Blur()
}

func (m ConnectionList) Init() tea.Cmd { return nil }

func (m ConnectionList) Update(msg tea.Msg) (ConnectionList, tea.Cmd) {
	var cmd tea.Cmd = nil

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ConnectionList) View() string {
	return getStyle(m.focus).Render(m.table.View())
}
