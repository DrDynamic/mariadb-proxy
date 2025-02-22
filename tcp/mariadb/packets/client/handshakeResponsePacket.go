package client

import (
	"fmt"
	"mschon/dbproxy/readers"
	"mschon/dbproxy/tcp/mariadb/packets"
)

const PacketTypeHandshakeResponsePacket packets.PacketType = "HandshakeResponsePacket"

type HandshakeResponsePacket struct {
	packets.Packet
	Client                   Client
	AuthPlugin               AuthPlugin
	ConnectionAttributeCount uint64
	ConnectionAttributes     map[string]string
}

type Client struct {
	Capabilities1    uint32
	Capabilities2    uint32
	Capabilities     uint64
	MaxPacketSize    uint32
	DefaultCollation uint8
	Username         string
	DefaultDatabase  string
}

type AuthPlugin struct {
	Data []byte
	Name string
}

func NewHandshakeResponsePacket(packet *packets.Packet) HandshakeResponsePacket {
	reader := readers.NewLeByteReader(packet.Body)

	capabilities1 := reader.PopUInt32()
	maxPkgSize := reader.PopUInt32()
	defaultCollation := reader.PopUInt8()
	reader.PopBytes(19) // reserved
	capabilities2 := reader.PopUInt32()

	capabilities := ((uint64(capabilities2) << 32) | uint64(capabilities1))

	username := reader.PopNullTerminatedString()

	var authData []byte
	if (capabilities & packets.PLUGIN_AUTH_LENENC_CLIENT_DATA) != 0 {
		authData = reader.PopLengthEncodedBytes()
	} else if (capabilities & packets.SECURE_CONNECTION) != 0 {
		len := reader.PopUInt8()
		authData = reader.PopBytes(int(len))
	} else {
		authData = reader.PopNullTerminatedBytes()
	}

	defaultDatabase := ""
	if (capabilities & packets.CONNECT_WITH_DB) != 0 {
		defaultDatabase = reader.PopNullTerminatedString()
	}

	authPluginName := ""
	if (capabilities & packets.PLUGIN_AUTH) != 0 {
		authPluginName = reader.PopNullTerminatedString()
	}

	var attrCount uint64 = 0
	attrs := make(map[string]string)
	if (capabilities & packets.CONNECT_ATTRS) != 0 {
		attrCount = reader.PopLengthEncodedInt()
		for reader.Length() > 0 {
			key := reader.PopLengthEncodedString()
			value := reader.PopLengthEncodedString()
			attrs[key] = value
		}
	}

	packet.Type = PacketTypeHandshakeResponsePacket

	return HandshakeResponsePacket{
		Packet: *packet,
		Client: Client{
			Capabilities1:    capabilities1,
			Capabilities2:    capabilities2,
			Capabilities:     capabilities,
			MaxPacketSize:    maxPkgSize,
			DefaultCollation: defaultCollation,
			Username:         username,
			DefaultDatabase:  defaultDatabase,
		},
		AuthPlugin: AuthPlugin{
			Data: authData,
			Name: authPluginName,
		},
		ConnectionAttributeCount: attrCount,
		ConnectionAttributes:     attrs,
	}
}

func (packet HandshakeResponsePacket) GetType() packets.PacketType {
	return packet.Type
}

func (packet HandshakeResponsePacket) GetSequence() byte {
	return packet.Sequence
}

func (packet HandshakeResponsePacket) String() string {
	return fmt.Sprintf("[HandshakeResponsePacket username=%s database=%s]", packet.Client.Username, packet.Client.DefaultDatabase)
}
