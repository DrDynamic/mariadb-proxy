package client

import (
	"mschon/dbproxy/readers"
	"mschon/dbproxy/tcp/mariadb/packets"
)

const PacketTypeComQuitPacket packets.PacketType = "ComQuitPacket"

type ComQuitPacket struct {
	packets.Packet
}

func NewComQuitPacket(packet *packets.Packet) ComQuitPacket {
	reader := readers.NewLeByteReader(packet.Body)
	reader.PopBytes(1) // header = 0x01
	packet.Type = PacketTypeComQueryPacket
	return ComQuitPacket{
		Packet: *packet,
	}
}

func (packet ComQuitPacket) GetType() packets.PacketType {
	return packet.Type
}

func (packet ComQuitPacket) GetSequence() byte {
	return packet.Sequence
}

func (packet ComQuitPacket) String() string {
	return "[ComQuitPacket]"
}
