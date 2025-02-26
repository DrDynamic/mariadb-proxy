package main

import (
	"fmt"
	"mschon/dbproxy/tcp"
	"mschon/dbproxy/tcp/mariadb"
	"mschon/dbproxy/tcp/mariadb/packets"
	"mschon/dbproxy/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config = GetConfig()

	//	p := tea.NewProgram(ui.New())
	var p *tea.Program = nil

	proxy := tcp.NewProxy()

	facrotyState := mariadb.NewPacketFactoryState()

	serverFactory := mariadb.NewServerPacketFactory(&facrotyState)
	clientFactory := mariadb.NewClientPacketFactory(&facrotyState)

	go onServerRec(&serverFactory, &proxy, p)
	go onClientRec(&clientFactory, &proxy, p)
	proxy.Listen(config.ProxyHost, config.ForwardHost)

	/*
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	*/
}

var test int = 0

func onServerRec(serverFactory *mariadb.ServerPacketFactory, proxy *tcp.Proxy, program *tea.Program) {

	//	builder := packets.NewPacketBuilder()
	for {
		msg := <-proxy.ServerRecChan
		/**
		builder.AddBytes(msg.Data)

		packet := builder.BuildPacket()
		for packet != nil {
			packet = builder.BuildPacket()
		}
		fmt.Printf("Server: %d bytes\n\r\n\r", len(builder.GetBuffer()))
		**/
		/**/
		serverFactory.AddBytes(msg.Data)
		var packet packets.BasePacket = serverFactory.CreatePacket()

		for packet != nil {
			fmt.Printf("[S][%02d] %s\r\n", packet.GetSequence(), packet.String())
			packet = serverFactory.CreatePacket()
		}

		//		fmt.Printf("ss %d  - %d\n\r", serverFactory.GetBufferSize(), serverFactory.GetState().ServerState)
		//		fmt.Printf("s %v\n\r", serverFactory.GetBuffer())

		for packet != nil && program != nil {
			program.Send(ui.PacketIoMessage{
				Direction:  msg.Direction,
				Connection: msg.Connection,
				Packet:     packet,
			})

			packet = serverFactory.CreatePacket()
		}
		/**/
	}

}

func onClientRec(clientFactory *mariadb.ClientPacketFactory, proxy *tcp.Proxy, program *tea.Program) {

	//builder := packets.NewPacketBuilder()
	for {
		msg := <-proxy.ClientRecChan
		/**
		builder.AddBytes(msg.Data)

		packet := builder.BuildPacket()
		for packet != nil {
			packet = builder.BuildPacket()
		}
		fmt.Printf("Client: %d bytes\n\r\n\r", len(builder.GetBuffer()))
		**/
		/**/
		clientFactory.AddBytes(msg.Data)
		var packet packets.BasePacket = clientFactory.CreatePacket()

		for packet != nil {
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
		/**/
	}

}
