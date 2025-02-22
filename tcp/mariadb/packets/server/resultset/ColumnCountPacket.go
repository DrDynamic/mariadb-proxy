package resultset

import (
	"fmt"
	"mschon/dbproxy/readers"
	"mschon/dbproxy/tcp/mariadb/packets"
)

const PacketTypeColumnCountPacket packets.PacketType = "ColumnCountPacket"

type ColumnCountPacket struct {
	packets.Packet
	ColumnCount     uint64
	MetadataFollows bool
}

func NewColumnCountPacket(packet *packets.Packet) ColumnCountPacket {
	reader := readers.NewLeByteReader(packet.Body)

	columnCount := reader.PopLengthEncodedInt()

	metadataFollows := false
	if reader.Length() > 0 && reader.PopBytes(1)[0] > 0 {
		metadataFollows = true
	}

	packet.Type = PacketTypeColumnCountPacket

	return ColumnCountPacket{
		Packet:          *packet,
		ColumnCount:     columnCount,
		MetadataFollows: metadataFollows,
	}
}

func (packet ColumnCountPacket) GetType() packets.PacketType {
	return packet.Type
}

func (packet ColumnCountPacket) GetSequence() byte {
	return packet.Sequence
}

func (packet ColumnCountPacket) String() string {
	return fmt.Sprintf("[ColumnCountPacket count=%d]", packet.ColumnCount)
}
