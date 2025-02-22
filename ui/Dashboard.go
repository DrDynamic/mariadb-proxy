package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Dashboard struct {
	packageList tea.Model
}

func NewDashboard() Dashboard {
	return Dashboard{
		packageList: NewPackageList(),
	}
}

func (m Dashboard) Init() tea.Cmd {
	return nil
}

func (m Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	m.packageList, cmd = m.packageList.Update(msg)

	return m, cmd
}

func (m Dashboard) View() string {
	return m.packageList.View()
}
