package mariadb

import (
	"mschon/dbproxy/tcp/mariadb/packets"
	"mschon/dbproxy/tcp/mariadb/packets/server/resultset"
)

type PacketFactory interface {
	AddBytes(data []byte)
	CreatePacket() packets.BasePacket
}

type PacketFactoryState struct {
	ServerState  serverPacketFactoryState
	ClientState  clientPacketFactoryState
	Capabilities uint64

	ResultsetInfo ResultsetInfo
}

type ResultsetInfo struct {
	IsBinaryProtocol  bool
	ColumnCount       int
	ColumnDefinitions []resultset.ColumnDefinitionPacket
}

type PacketIdentifier func(packet *packets.Packet) packets.BasePacket

func NewPacketFactoryState() PacketFactoryState {
	return PacketFactoryState{
		ResultsetInfo: ResultsetInfo{
			IsBinaryProtocol:  false,
			ColumnCount:       0,
			ColumnDefinitions: make([]resultset.ColumnDefinitionPacket, 0),
		},
	}
}

func identifyPacket(packet *packets.Packet, identifiers []PacketIdentifier) packets.BasePacket {
	var result packets.BasePacket = nil
	for _, identifier := range identifiers {
		result = identifier(packet)

		switch result.(type) {
		case *packets.Packet:
			break
		default:
			return result
		}
	}
	return result
}
