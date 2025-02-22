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

	proxy := tcp.NewProxy()

	go onRec(&proxy, p)
	go proxy.Listen(config.ProxyHost, config.ForwardHost)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	/*
	   config = GetConfig()

	   listener, err := net.Listen("tcp", config.ProxyHost)

	   	if err != nil {
	   		fmt.Println("Listen error: ", err)
	   	}

	   conn, err := listener.Accept()

	   	if err != nil {
	   		fmt.Println("Accept error: ", err)
	   	}

	   fwd, err := net.Dial("tcp", config.ForwardHost)

	   	if err != nil {
	   		fmt.Println("Dial error: ", err)
	   	}

	   handler := mariadb.NewConnectionHandler()

	   go proxyHandler("local 2 fwd", conn, fwd, handler.ClientPacket)
	   proxyHandler("fwd 2 local", fwd, conn, handler.ServerPacket)
	*/
}

func onRec(proxy *tcp.Proxy, program *tea.Program) {

	facrotyState := mariadb.NewPacketFactoryState()

	serverFactory := mariadb.NewServerPacketFactory(&facrotyState)
	clientFactory := mariadb.NewClientPacketFactory(&facrotyState)

	for {
		msg := <-proxy.RecChan

		var factory mariadb.PacketFactory

		switch msg.Direction {
		case tcp.ProxyRecDirection_fromServer:
			factory = &serverFactory
		case tcp.ProxyRecDirection_toServer:
			factory = &clientFactory

		}

		factory.AddBytes(msg.Data)
		var packet packets.BasePacket = factory.CreatePacket()

		for packet != nil {
			program.Send(ui.PacketIoMessage{
				Direction:  msg.Direction,
				Connection: msg.Connection,
				Packet:     packet,
			})

			packet = factory.CreatePacket()
		}
	}

}
