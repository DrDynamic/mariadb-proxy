package resultset

import (
	"fmt"
	"mschon/dbproxy/readers"
	"mschon/dbproxy/tcp/mariadb/packets"
)

const PacketTypeResultsetRowPacket packets.PacketType = "ResultsetRowPacket"

type ResultsetRowPacket struct {
	packets.Packet
	Columns []string
}

func NewResultsetRowPacket(packet *packets.Packet) ResultsetRowPacket {
	reader := readers.NewLeByteReader(packet.Body)

	columns := make([]string, 0)

	for reader.Length() > 0 {
		columns = append(columns, reader.PopLengthEncodedString())
	}

	packet.Type = PacketTypeColumnCountPacket

	return ResultsetRowPacket{
		Packet:  *packet,
		Columns: columns,
	}
}

func (packet ResultsetRowPacket) GetType() packets.PacketType {
	return packet.Type
}

func (packet ResultsetRowPacket) GetSequence() byte {
	return packet.Sequence
}

func (packet ResultsetRowPacket) String() string {
	return fmt.Sprintf("[ResultsetRowPacket count=%d]", len(packet.Columns))
}
