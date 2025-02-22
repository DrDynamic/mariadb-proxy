package server

import (
	"fmt"
	"mschon/dbproxy/readers"
	"mschon/dbproxy/tcp/mariadb/packets"
)

const PacketTypeLocalInlinePacket packets.PacketType = "LocalInlinePacket"

type LocalInlinePacket struct {
	packets.Packet
	Filename string
}

func NewLocalInlinePacket(packet *packets.Packet) LocalInlinePacket {
	reader := readers.NewLeByteReader(packet.Body)
	reader.PopBytes(1) // header = 0xFB

	packet.Type = PacketTypeLocalInlinePacket

	return LocalInlinePacket{
		Packet:   *packet,
		Filename: reader.PopEofEncodedString(),
	}
}

func (packet LocalInlinePacket) GetType() packets.PacketType {
	return packet.Type
}

func (packet LocalInlinePacket) GetSequence() byte {
	return packet.Sequence
}

func (packet LocalInlinePacket) String() string {
	return fmt.Sprintf("[LocalInlinePacket sql=%s]", packet.Filename)
}
