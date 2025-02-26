package ui

import (
	"mschon/dbproxy/tcp"
	"mschon/dbproxy/tcp/mariadb/packets"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type PackageIoDirection int

const (
	DirectionToServer PackageIoDirection = iota
	DirectionFromServer
)

type PacketIoMessage struct {
	Direction  tcp.ProxyRecDirection
	Connection tcp.ProxyConnection
	Packet     packets.BasePacket
}

type PackageList struct {
	table table.Model
}

var columns = []table.Column{
	{Title: "<>", Width: 4},
	{Title: "Seq", Width: 4},
	{Title: "Package", Width: 100},
}

func NewPackageList() PackageList {
	return PackageList{
		table: table.New(
			table.WithColumns(columns),
		),
	}
}

func (m PackageList) Init() tea.Cmd {
	return nil
}

func (m PackageList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case PacketIoMessage:

		rows := m.table.Rows()

		rows = append(rows, table.Row{msg.Direction.String(), strconv.Itoa(int(msg.Packet.GetSequence())), msg.Packet.String()})
		m.table.SetRows(rows)

		//		fmt.Printf("[%d] %s\r\n", msg.Packet.GetSequence(), msg.Packet.String())

	default:
		m.table, cmd = m.table.Update(msg)
	}
	return m, cmd
}

func (m PackageList) View() string {
	return "" //baseStyle.Render(m.table.View()) + "\n" + m.table.HelpView()
}
