package mariadb

import (
	"mschon/dbproxy/tcp/mariadb/packets"
	"mschon/dbproxy/tcp/mariadb/packets/server"
	"mschon/dbproxy/tcp/mariadb/packets/server/resultset"
)

type serverPacketFactoryState int

const (
	ServerStateNone serverPacketFactoryState = iota

	ServerStateInit
	ServerStateInitFinish
	ServerStateQueryResponse
	ServerStateResultsetColumnDefinition
	ServerStateResultsetEof
	ServerStateResultsetRow
)

type ServerPacketFactory struct {
	state   *PacketFactoryState
	builder packets.PacketBuilder

	RecPackets int
}

func NewServerPacketFactory(state *PacketFactoryState) ServerPacketFactory {
	state.ServerState = ServerStateInit

	builder := packets.NewPacketBuilder()
	builder.DebugPrint = false

	return ServerPacketFactory{
		state:   state,
		builder: builder,

		RecPackets: 0,
	}
}

func (factory *ServerPacketFactory) GetState() PacketFactoryState {
	return *factory.state
}

func (factory *ServerPacketFactory) GetBufferSize() int {
	return len(factory.builder.GetBuffer())
}
func (factory *ServerPacketFactory) GetBuffer() []byte {
	return factory.builder.GetBuffer()
}

func (factory *ServerPacketFactory) AddBytes(data []byte) {
	factory.builder.AddBytes(data)
}

func (factory *ServerPacketFactory) CreatePacket() packets.BasePacket {
	packet := factory.builder.BuildPacket()
	if packet == nil {
		return nil
	}

	factory.RecPackets++
	return factory.convertPacket(packet)
}
func (factory *ServerPacketFactory) convertPacket(packet *packets.Packet) packets.BasePacket {

	switch factory.state.ServerState {
	case ServerStateInit:
		factory.state.ServerState = ServerStateInitFinish
		handshake := server.NewInitialHandshakePacket(packet)

		factory.state.Capabilities = handshake.Server.Capabilities

		return handshake
	case ServerStateInitFinish:
		factory.state.ServerState = ServerStateNone
		return identifyPacket(packet, []PacketIdentifier{
			factory.identifyOkPacket,
			factory.identifyErrPacket,
		})
	case ServerStateQueryResponse:
		p := identifyPacket(packet, []PacketIdentifier{
			factory.identifyOkPacket,
			factory.identifyErrPacket,
			factory.identifyLocalInlinePacket,
			factory.identifyResultSet,
		})

		if _, ok := p.(resultset.ColumnCountPacket); ok {
			cc := p.(resultset.ColumnCountPacket)

			factory.state.ResultsetInfo.ColumnCount = int(cc.ColumnCount)
			factory.state.ResultsetInfo.ColumnDefinitions = make([]resultset.ColumnDefinitionPacket, 0)

			if (factory.state.Capabilities&packets.MARIADB_CLIENT_CACHE_METADATA) == 0 ||
				cc.MetadataFollows {
				factory.state.ServerState = ServerStateResultsetColumnDefinition
			} else {
				factory.state.ServerState = ServerStateResultsetRow
			}
		} else {
			factory.state.ServerState = ServerStateNone
		}

		return p
	case ServerStateResultsetColumnDefinition:

		var definition packets.BasePacket = nil
		if factory.state.ResultsetInfo.ColumnCount > len(factory.state.ResultsetInfo.ColumnDefinitions) {
			definition = resultset.NewColumnDefinitionPacket(packet, factory.state.Capabilities)
			factory.state.ResultsetInfo.ColumnDefinitions = append(factory.state.ResultsetInfo.ColumnDefinitions, definition.(resultset.ColumnDefinitionPacket))
		} else {
			definition = factory.convertPacket(packet)
		}

		if factory.state.ResultsetInfo.ColumnCount <= len(factory.state.ResultsetInfo.ColumnDefinitions) {
			if (factory.state.Capabilities & packets.CLIENT_DEPRECATE_EOF) == 0 {
				factory.state.ServerState = ServerStateResultsetEof
			} else {
				factory.state.ServerState = ServerStateResultsetRow
			}
		}

		return definition

	case ServerStateResultsetEof:
		return factory.identifyEofPacket(packet)
	case ServerStateResultsetRow:
		p := identifyPacket(packet, []PacketIdentifier{
			factory.identifyEofPacket,
			factory.identifyErrPacket,
			factory.identifyOkPacket,
			factory.identifyResultSetRow,
		})

		if _, ok := p.(resultset.ResultsetRowPacket); !ok {
			factory.state.ServerState = ServerStateNone
		}

		return p
	case ServerStateNone:
		return packet
	}

	return packet
}

func (factory *ServerPacketFactory) identifyEofPacket(packet *packets.Packet) packets.BasePacket {
	if (factory.state.Capabilities&packets.CLIENT_DEPRECATE_EOF) == 0 &&
		packet.Body[0] == 0xFE {
		return server.NewEofPacket(packet)
	}
	return packet
}

func (factory *ServerPacketFactory) identifyErrPacket(packet *packets.Packet) packets.BasePacket {
	if packet.Body[0] == 0xFF {
		return server.NewErrPacket(packet)
	}
	return packet
}

func (factory *ServerPacketFactory) identifyOkPacket(packet *packets.Packet) packets.BasePacket {
	if packet.Body[0] == 0x00 {
		return server.NewOkPacket(packet)
	}

	if (factory.state.Capabilities&packets.CLIENT_DEPRECATE_EOF) != 0 &&
		packet.Body[0] == 0xFE {
		return server.NewOkPacket(packet)
	}

	return packet
}

func (factory *ServerPacketFactory) identifyLocalInlinePacket(packet *packets.Packet) packets.BasePacket {
	if packet.Body[0] == 0xFB {
		return server.NewLocalInlinePacket(packet)
	}
	return packet
}

func (factory *ServerPacketFactory) identifyResultSet(packet *packets.Packet) packets.BasePacket {
	if packet.Length == 1 || packet.Length == 2 {
		return resultset.NewColumnCountPacket(packet)
	}
	return packet
}

func (factory *ServerPacketFactory) identifyResultSetRow(packet *packets.Packet) packets.BasePacket {

	if packet.Body[0] == 0x00 ||
		packet.Body[0] == 0xFE ||
		packet.Body[0] == 0xFF {
		return packet
	}

	return resultset.NewResultsetRowPacket(packet)
}
