package mariadb

import (
	"fmt"
	"mschon/dbproxy/tcp/mariadb/packets"
	"mschon/dbproxy/tcp/mariadb/packets/client"
	"mschon/dbproxy/tcp/mariadb/packets/server"
)

type ConnectionHandler struct {
	serverPacketIndex int
	clientPacketIndex int
}

func NewConnectionHandler() *ConnectionHandler {
	return &ConnectionHandler{
		serverPacketIndex: 0,
		clientPacketIndex: 0,
	}
}

func (handler *ConnectionHandler) ServerPacket(packet *packets.Packet) {
	if handler.serverPacketIndex == 0 {
		fmt.Println("[SRV]", packet.Sequence, server.NewInitialHandshakePacket(packet))
	} else if handler.serverPacketIndex == 1 {
		fmt.Println("[SRV]", packet.Sequence, identifyServerPacket(packet))
	} else {
		fmt.Println("[SRV]", packet.Sequence, identifyServerPacket(packet))
	}
	handler.serverPacketIndex++
}

func identifyServerPacket(packet *packets.Packet) string {
	switch packet.Body[0] {
	case 0x00:
		return server.NewOkPacket(packet).String()
	case 0xFE:
		return server.NewOkPacket(packet).String()
	case 0xFF:
		return server.NewErrPacket(packet).String()
	default:
		return packet.String()
	}
}

func (handler *ConnectionHandler) ClientPacket(packet *packets.Packet) {
	if handler.clientPacketIndex == 0 {
		fmt.Println("[CLN]", packet.Sequence, client.NewHandshakeResponsePacket(packet))
	} else if handler.clientPacketIndex == 1 {
		fmt.Println("[CLN]", packet.Sequence, identifyClientPacket(packet))
	} else {
		fmt.Println("[CLN]", packet.Sequence, identifyClientPacket(packet))
	}
	handler.clientPacketIndex++
}

func identifyClientPacket(packet *packets.Packet) string {
	switch packet.Body[0] {
	//	case 0x03:
	//		return client.NewComQueryPacket(packet).String()
	default:
		return packet.String()
	}
}
