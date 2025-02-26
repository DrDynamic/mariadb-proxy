package server

import (
	"fmt"
	"math"
	"mschon/dbproxy/readers"
	"mschon/dbproxy/tcp/mariadb/packets"
)

const PacketTypeInitialHandshakePacket packets.PacketType = "InitialHandshakePacket"

type InitialHandshakePacket struct {
	packets.Packet
	ProtocolVersion uint8
	ConnectionId    uint32
	Server          Server
	AuthPlugin      AuthPlugin
}

type AuthPlugin struct {
	Data1 []byte
	Data2 []byte
	Name  string
}

type Server struct {
	Version          string
	Capabilities1    uint16
	Capabilities2    uint16
	Capabilities3    uint32
	Capabilities     uint64
	DefaultCollation uint8
	Status           uint16 // das hier? https://mariadb-corporation.github.io/mariadb-connector-python/constants.html#status
}

func NewInitialHandshakePacket(packet *packets.Packet) InitialHandshakePacket {
	reader := readers.NewLeByteReader(packet.Body)

	protocolVersion := reader.PopUInt8()
	serverVersion := reader.PopNullTerminatedString()

	connectionId := reader.PopUInt32()

	authPluginData1 := reader.PopBytes(8)
	reader.PopBytes(1) // reserved
	serverCapabilities1 := reader.PopUInt16()
	defaultCollation := reader.PopUInt8()
	statusFlags := reader.PopUInt16()
	serverCapabilities2 := reader.PopUInt16()
	pluginLength := reader.PopUInt8()
	reader.PopBytes(6) // filler
	serverCapabilities3 := reader.PopUInt32()

	serverCapabilities := ((uint64(serverCapabilities3) << 32) | (uint64(serverCapabilities2) << 16) | uint64(serverCapabilities1))

	// if capabilitiers.CLIENT_SECURE_CONNECTION
	size := int(math.Max(12, float64(pluginLength)-9))

	authPluginData2 := reader.PopBytes(size)
	reader.PopBytes(1) // reserved

	// if capabilities.PLUGIN_AUTH
	authPluginName := reader.PopNullTerminatedString()

	packet.Type = PacketTypeInitialHandshakePacket

	return InitialHandshakePacket{
		Packet:          *packet,
		ProtocolVersion: protocolVersion,
		ConnectionId:    connectionId,
		Server: Server{
			Version:          serverVersion,
			Capabilities1:    serverCapabilities1,
			Capabilities2:    serverCapabilities2,
			Capabilities3:    serverCapabilities3,
			Capabilities:     serverCapabilities,
			DefaultCollation: defaultCollation,
			Status:           statusFlags,
		},
		AuthPlugin: AuthPlugin{
			Data1: authPluginData1,
			Data2: authPluginData2,
			Name:  authPluginName,
		},
	}
}

func (packet InitialHandshakePacket) GetType() packets.PacketType {
	return packet.Type
}

func (packet InitialHandshakePacket) GetSequence() byte {
	return packet.Sequence
}

func (packet InitialHandshakePacket) String() string {
	return fmt.Sprintf("[InitialHandshakePacket protocol=%d server=%s connection=%d]", packet.ProtocolVersion, packet.Server.Version, packet.ConnectionId)
}
