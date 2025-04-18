package main

import (
	"fmt"
	"mschon/dbproxy/tcp"
	"mschon/dbproxy/tcp/mariadb"
	"mschon/dbproxy/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config = GetConfig()

	proxy := tcp.NewProxy()
	connectionManager := mariadb.NewConnectionManager(&proxy)

	p := tea.NewProgram(ui.New(&connectionManager))

	startProxy(&connectionManager, config, p)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func startProxy(connectionManager *mariadb.MariadbConnectionManager, config *Config, program *tea.Program) {

	connectionManager.AddOnConnectionChangeListener(func(connection mariadb.MariadbConnection) {
		program.Send(ui.UpdateConnectionsMsg{})
	})

	connectionManager.AddOnNewConnectionListener(func(connection mariadb.MariadbConnection) {
		program.Send(ui.UpdateConnectionsMsg{})
	})

	go connectionManager.Listen(config.ProxyHost, config.ForwardHost)
}
