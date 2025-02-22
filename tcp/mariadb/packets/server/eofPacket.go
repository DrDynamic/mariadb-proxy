package server

import (
	"fmt"
	"mschon/dbproxy/readers"
	"mschon/dbproxy/tcp/mariadb/packets"
)

const PacketTypeEofPacket packets.PacketType = "EofPacket"

type EofPacket struct {
	packets.Packet
	WarningCount uint16
	ServerStatus uint16
}

func NewEofPacket(packet *packets.Packet) EofPacket {
	reader := readers.NewLeByteReader(packet.Body)

	reader.PopUInt8() // always 0xFE

	serverStatus := reader.PopUInt16()
	warningCount := reader.PopUInt16()

	packet.Type = PacketTypeEofPacket

	return EofPacket{
		Packet:       *packet,
		WarningCount: warningCount,
		ServerStatus: serverStatus,
	}
}

func (packet EofPacket) GetType() packets.PacketType {
	return packet.Type
}

func (packet EofPacket) GetSequence() byte {
	return packet.Sequence
}

func (packet EofPacket) String() string {
	return fmt.Sprintf("[EofPacket warnings=%d status=%d]", packet.WarningCount, packet.ServerStatus)
}
