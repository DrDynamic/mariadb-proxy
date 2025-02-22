package client

import (
	"fmt"
	"mschon/dbproxy/readers"
	"mschon/dbproxy/tcp/mariadb/packets"
)

const PacketTypeComQueryPacket packets.PacketType = "ComQueryPacket"

type ComQueryPacket struct {
	packets.Packet
	SqlStatement string
}

func NewComQueryPacket(packet *packets.Packet) ComQueryPacket {
	reader := readers.NewLeByteReader(packet.Body)
	reader.PopBytes(1) // header = 0x03
	packet.Type = PacketTypeComQueryPacket
	return ComQueryPacket{
		Packet:       *packet,
		SqlStatement: reader.PopEofEncodedString(),
	}
}

func (packet ComQueryPacket) GetType() packets.PacketType {
	return packet.Type
}

func (packet ComQueryPacket) GetSequence() byte {
	return packet.Sequence
}

func (packet ComQueryPacket) String() string {
	return fmt.Sprintf("[ComQueryPacket sql=%s]", packet.SqlStatement)
}
