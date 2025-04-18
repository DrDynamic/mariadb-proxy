package ui

import (
	"fmt"
	"mschon/dbproxy/tcp/mariadb"
	"mschon/dbproxy/util"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UpdateConnectionsMsg struct {
}

type ConnectionList struct {
	focus             bool
	table             table.Model
	connectionManager *mariadb.MariadbConnectionManager
}

func makeRow(connection *mariadb.MariadbConnection) table.Row {
	proxyConnection := connection.GetProxyConnection()

	return table.Row{fmt.Sprintf("%p", proxyConnection), string(proxyConnection.Status)}
}

func NewConnectionList(connectionManager *mariadb.MariadbConnectionManager) ConnectionList {
	var columns = []table.Column{
		{Title: "id", Width: 16},
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

	rows := util.MapSlice(connectionManager.Connections, makeRow)

	return ConnectionList{
		focus: false,
		table: table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(false),
			table.WithStyles(styles),
		),
		connectionManager: connectionManager,
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

	switch msg.(type) {
	case UpdateConnectionsMsg:
		rows := util.MapSlice(m.connectionManager.Connections, makeRow)
		m.table.SetRows(rows)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ConnectionList) View() string {
	return getStyle(m.focus).Render(m.table.View())
}
