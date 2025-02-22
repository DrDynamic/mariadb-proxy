package resultset

import (
	"fmt"
	"mschon/dbproxy/readers"
	"mschon/dbproxy/tcp/mariadb/packets"
)

const PacketTypeColumnDefinitionPacket packets.PacketType = "ColumnDefinitionPacket"

type ColumnDefinitionPacket struct {
	packets.Packet
	Schema           string
	TableAlias       string
	Table            string
	ColumnAlias      string
	Column           string
	ExtendedMeta     string
	FixedFieldLength uint64
	CharacterSet     uint16
	MaxColumnSize    uint32
	FieldTypes       uint8
	FieldDetailFlag  uint16
	Decimanls        uint8
}

func NewColumnDefinitionPacket(packet *packets.Packet, capabilities uint64) ColumnDefinitionPacket {
	reader := readers.NewLeByteReader(packet.Body)
	reader.PopBytes(4) // always 3 "def"

	schema := reader.PopLengthEncodedString()
	tableAlias := reader.PopLengthEncodedString()
	table := reader.PopLengthEncodedString()
	columnAlias := reader.PopLengthEncodedString()
	column := reader.PopLengthEncodedString()

	extendedMeta := ""
	if (capabilities & packets.MARIADB_CLIENT_EXTENDED_METADATA) != 0 {
		extendedMeta = reader.PopLengthEncodedString()
	}

	fixedFieldLength := reader.PopLengthEncodedInt()
	charSet := reader.PopUInt16()
	maxColumnSize := reader.PopUInt32()
	fieldTypes := reader.PopUInt8()
	fieldDetails := reader.PopUInt16()
	decimals := reader.PopUInt8()

	packet.Type = PacketTypeColumnDefinitionPacket

	return ColumnDefinitionPacket{
		Packet:           *packet,
		Schema:           schema,
		TableAlias:       tableAlias,
		Table:            table,
		ColumnAlias:      columnAlias,
		Column:           column,
		ExtendedMeta:     extendedMeta,
		FixedFieldLength: fixedFieldLength,
		CharacterSet:     charSet,
		MaxColumnSize:    maxColumnSize,
		FieldTypes:       fieldTypes,
		FieldDetailFlag:  fieldDetails,
		Decimanls:        decimals,
	}
}

func (packet ColumnDefinitionPacket) GetType() packets.PacketType {
	return packet.Type
}

func (packet ColumnDefinitionPacket) GetSequence() byte {
	return packet.Sequence
}

func (packet ColumnDefinitionPacket) String() string {
	//	return fmt.Sprintf("[ColumnDefinitionPacket Schema=%s Table=%s Column=%s]", packet.Schema, packet.TableAlias, packet.ColumnAlias)
	return fmt.Sprintf("[ColumnDefinitionPacket Schema=%s Table=%s Column=%s]", packet.Schema, packet.Table, packet.Column)
}
