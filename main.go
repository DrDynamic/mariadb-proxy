package main

import (
	"fmt"
	"mschon/dbproxy/tcp"
	"mschon/dbproxy/tcp/mariadb"
	"mschon/dbproxy/tcp/mariadb/packets"
	"mschon/dbproxy/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config = GetConfig()

	p := tea.NewProgram(ui.New())
	//var p *tea.Program = nil

	startProxy(config, p)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func startProxy(config *Config, program *tea.Program) {
	proxy := tcp.NewProxy()

	facrotyState := mariadb.NewPacketFactoryState()

	serverFactory := mariadb.NewServerPacketFactory(&facrotyState)
	clientFactory := mariadb.NewClientPacketFactory(&facrotyState)

	go onServerRec(&serverFactory, &proxy, program)
	go onClientRec(&clientFactory, &proxy, program)
	go proxy.Listen(config.ProxyHost, config.ForwardHost)
}

func onServerRec(serverFactory *mariadb.ServerPacketFactory, proxy *tcp.Proxy, program *tea.Program) {

	for {
		msg := <-proxy.ServerRecChan

		serverFactory.AddBytes(msg.Data)
		var packet packets.BasePacket = serverFactory.CreatePacket()

		for packet != nil && program == nil {
			fmt.Printf("[S][%02d] %s\r\n", packet.GetSequence(), packet.String())
			packet = serverFactory.CreatePacket()
		}

		for packet != nil && program != nil {
			program.Send(ui.PacketIoMessage{
				Direction:  msg.Direction,
				Connection: msg.Connection,
				Packet:     packet,
			})

			packet = serverFactory.CreatePacket()
		}
	}

}

func onClientRec(clientFactory *mariadb.ClientPacketFactory, proxy *tcp.Proxy, program *tea.Program) {

	for {
		msg := <-proxy.ClientRecChan

		clientFactory.AddBytes(msg.Data)
		var packet packets.BasePacket = clientFactory.CreatePacket()

		for packet != nil && program == nil {
			fmt.Printf("[C][%02d] %s\r\n", packet.GetSequence(), packet.String())
			packet = clientFactory.CreatePacket()
		}

		for packet != nil && program != nil {
			program.Send(ui.PacketIoMessage{
				Direction:  msg.Direction,
				Connection: msg.Connection,
				Packet:     packet,
			})

			packet = clientFactory.CreatePacket()
		}
	}

}
