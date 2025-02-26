package mariadb

import (
	"mschon/dbproxy/tcp/mariadb/packets"
	"mschon/dbproxy/tcp/mariadb/packets/client"
)

type clientPacketFactoryState int

const (
	ClientStateNone clientPacketFactoryState = iota

	ClientStateHandshakeResponse
	ClientStateConected
)

type ClientPacketFactory struct {
	state   *PacketFactoryState
	builder packets.PacketBuilder
}

func NewClientPacketFactory(state *PacketFactoryState) ClientPacketFactory {
	state.ClientState = ClientStateHandshakeResponse
	return ClientPacketFactory{
		state:   state,
		builder: packets.NewPacketBuilder(),
	}
}

func (factory *ClientPacketFactory) GetBufferSize() int {
	return len(factory.builder.GetBuffer())
}

func (factory *ClientPacketFactory) AddBytes(data []byte) {
	factory.builder.AddBytes(data)
}

func (factory *ClientPacketFactory) CreatePacket() packets.BasePacket {
	packet := factory.builder.BuildPacket()
	if packet == nil {
		return nil
	}

	switch factory.state.ClientState {
	case ClientStateHandshakeResponse:
		factory.state.ClientState = ClientStateConected
		response := client.NewHandshakeResponsePacket(packet)
		factory.state.Capabilities = response.Client.Capabilities
		return response
	case ClientStateConected:

		return identifyPacket(packet, []PacketIdentifier{
			factory.identifyClientQueryPacket,
			factory.identifyClientQuitPacket,
		})
	}

	return packet
}

func (factory *ClientPacketFactory) identifyClientQueryPacket(packet *packets.Packet) packets.BasePacket {
	id := packet.Body[0]

	if id == 0x03 {
		factory.state.ServerState = ServerStateQueryResponse
		factory.state.ResultsetInfo.IsBinaryProtocol = false
		return client.NewComQueryPacket(packet)
	}

	return packet
}

func (factory *ClientPacketFactory) identifyClientQuitPacket(packet *packets.Packet) packets.BasePacket {
	id := packet.Body[0]

	if id == 0x01 {
		factory.state.ServerState = ServerStateNone
		return client.NewComQuitPacket(packet)
	}

	return packet
}
