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

type clientPacketFactory struct {
	state   *PacketFactoryState
	builder packets.PacketBuilder
}

func NewClientPacketFactory(state *PacketFactoryState) clientPacketFactory {
	state.ClientState = ClientStateHandshakeResponse
	return clientPacketFactory{
		state:   state,
		builder: *packets.NewPacketBuilder(),
	}
}

func (factory *clientPacketFactory) AddBytes(data []byte) {
	factory.builder.AddBytes(data)
}

func (factory *clientPacketFactory) CreatePacket() packets.BasePacket {
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
		})
	}

	return packet
}

func (factory *clientPacketFactory) identifyClientQueryPacket(packet *packets.Packet) packets.BasePacket {
	id := packet.Body[0]

	if id == 0x03 {
		factory.state.ServerState = ServerStateQueryResponse
		factory.state.ResultsetInfo.IsBinaryProtocol = false
		return client.NewComQueryPacket(packet)
	}

	return packet
}
