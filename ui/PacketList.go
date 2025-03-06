package ui

import (
	"mschon/dbproxy/tcp"
	"mschon/dbproxy/tcp/mariadb/packets"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PacketIoMessage struct {
	Direction  tcp.ProxyRecDirection
	Connection tcp.ProxyConnection
	Packet     packets.BasePacket
}

type PacketList struct {
	focus bool
	table table.Model
}

var columns = []table.Column{
	{Title: "<>", Width: 4},
	{Title: "Seq", Width: 4},
	{Title: "Packet", Width: 100},
}

func NewPacketList() PacketList {
	styles := table.DefaultStyles()
	styles.Selected = styles.Selected.
		Background(lipgloss.Color("240")).
		UnsetForeground()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true)

	return PacketList{
		focus: false,
		table: table.New(
			table.WithColumns(columns),
			table.WithRows([]table.Row{
				{">", "0", "Demo 1"},
				{"<", "1", "Demo 2"},
				{">", "2", "Demo 3"},
				{"<", "3", "Demo 4"},
				{">", "4", "Demo 5"},
				{"<", "5", "Demo 6"},
			}),
			table.WithStyles(styles),
		),
	}
}

func (m *PacketList) GetFocus() bool {
	return m.focus
}

func (m *PacketList) Focus() {
	m.focus = true
	m.table.Focus()
}

func (m *PacketList) Blur() {
	m.focus = false
	m.table.Blur()
}

func (m PacketList) Init() tea.Cmd { return nil }

func (m PacketList) Update(msg tea.Msg) (PacketList, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case PacketIoMessage:
		rows := m.table.Rows()
		rows = append(rows, table.Row{msg.Direction.String(), strconv.Itoa(int(msg.Packet.GetSequence())), msg.Packet.String()})
		m.table.SetRows(rows)

	default:
		m.table, cmd = m.table.Update(msg)
	}
	return m, cmd
}

func (m PacketList) View() string {
	return getStyle(m.focus).Render(m.table.View()) // + "\n" + m.table.HelpView()
}
