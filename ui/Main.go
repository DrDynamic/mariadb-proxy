package ui

import (
	"fmt"
	"mschon/dbproxy/tcp/mariadb"

	tea "github.com/charmbracelet/bubbletea"
)

type FocusMsg struct{}
type BlurMsg struct{}

type Focusable interface {
	GetFocus() bool
	Focus()
	Blur()
}

type main struct {
	view tea.Model
}

func New(connectionManager *mariadb.MariadbConnectionManager) main {
	return main{
		view: NewDashboard(connectionManager),
	}
}

func (m main) Init() tea.Cmd { return nil }

func (m main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			fmt.Println("[Quit message]", msg)
			return m, tea.Quit
		}

	}

	var cmd tea.Cmd = nil
	m.view, cmd = m.view.Update(msg)
	return m, cmd
}

func (m main) View() string {
	return m.view.View()
}
