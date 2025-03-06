package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Dashboard struct {
	connectionList *ConnectionList
	packetList     *PacketList

	currentFocus int
	focusList    []Focusable
}

func NewDashboard() Dashboard {
	connectionList := NewConnectionList()
	packetList := NewPacketList()

	connectionList.Focus()

	return Dashboard{
		connectionList: &connectionList,
		packetList:     &packetList,

		currentFocus: 0,
		focusList: []Focusable{
			&connectionList,
			&packetList,
		},
	}
}

func (m Dashboard) Init() tea.Cmd {
	return nil
}

type FCS struct {
	focus bool
}

func (FCS) Blur()  {}
func (FCS) Focus() {}

func (m Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.currentFocus >= len(m.focusList)-1 {
				m.currentFocus = 0
			} else {
				m.currentFocus += 1
			}
			m.updateFocus()
		case "shift+tab":
			if m.currentFocus <= 0 {
				m.currentFocus = len(m.focusList) - 1
			} else {
				m.currentFocus -= 1
			}
			m.updateFocus()
		}
	}

	*m.connectionList, _ = m.connectionList.Update(msg)
	*m.packetList, _ = m.packetList.Update(msg)

	return m, cmd
}

func (m *Dashboard) updateFocus() {
	for index, model := range m.focusList {
		model.Blur()
		if index == m.currentFocus {
			model.Focus()
		}
	}
}

func (m Dashboard) View() string {
	return lipgloss.JoinHorizontal(0, m.connectionList.View(), m.packetList.View())
}
