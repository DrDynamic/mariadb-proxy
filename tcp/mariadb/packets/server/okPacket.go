package server

import (
	"fmt"
	"mschon/dbproxy/readers"
	"mschon/dbproxy/tcp/mariadb/packets"
)

const PacketTypeOkPacket packets.PacketType = "OkPacket"

type OkPacket struct {
	packets.Packet
	DeprecateEOF bool
	AffectedRows int
	LastInsertId int
	ServerStatus uint16
	WarningCount uint16
	Info         []byte
}

func NewOkPacket(packet *packets.Packet) *OkPacket {
	reader := readers.NewLeByteReader(packet.Body)

	head := reader.PopUInt8()
	affectedRows := reader.PopLengthEncodedInt()
	lastInsertId := reader.PopLengthEncodedInt()
	serverStatus := reader.PopUInt16()
	warningCount := reader.PopUInt16()

	packet.Type = PacketTypeOkPacket

	return &OkPacket{
		Packet:       *packet,
		DeprecateEOF: head == 0xFE,
		AffectedRows: int(affectedRows),
		LastInsertId: int(lastInsertId),
		ServerStatus: serverStatus,
		WarningCount: warningCount,
	}
}

func (packet OkPacket) GetType() packets.PacketType {
	return packet.Type
}

func (packet OkPacket) GetSequence() byte {
	return packet.Sequence
}

func (packet *OkPacket) String() string {
	return fmt.Sprintf("[OkPacket rows=%d lastId=%d warnings=%d]",
		packet.AffectedRows,
		packet.LastInsertId,
		packet.WarningCount)
}
