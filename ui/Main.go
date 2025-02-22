package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type main struct {
	view tea.Model
}

func New() main {
	return main{
		view: NewDashboard(),
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
