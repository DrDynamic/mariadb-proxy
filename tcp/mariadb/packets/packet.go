package packets

import (
	"encoding/binary"
	"fmt"
)

type PacketType string

type BasePacket interface {
	GetType() PacketType
	GetSequence() byte

	String() string
}

type Packet struct {
	Type PacketType

	Length   uint32
	Sequence byte
	Body     []byte

	IsComplete bool
}

const PacketTypeGeneral PacketType = "General"

func PacketFromBytes(data []byte) *Packet {
	length := binary.LittleEndian.Uint32([]byte{data[0], data[1], data[2], 0x00})

	var body []byte = nil
	var isComplete = false

	if len(data) < int(length)+4 {
		body = data[4:]
		isComplete = false
	} else {
		body = data[4 : 4+length]
		isComplete = true
	}

	return &Packet{
		Type:       PacketTypeGeneral,
		Length:     length,
		Sequence:   data[3],
		Body:       body,
		IsComplete: isComplete,
	}
}

func (packet *Packet) ExtendIncompletePacket(data []byte) {
	newLength := len(packet.Body) + len(data)

	if newLength < int(packet.Length)+4 {
		packet.Body = append(packet.Body, data...)
		packet.IsComplete = false
	} else {
		missing := int(packet.Length) - len(packet.Body)
		packet.Body = append(packet.Body, data[:missing]...)
		packet.IsComplete = true
	}
}

func (packet Packet) GetType() PacketType {
	return packet.Type
}

func (packet Packet) GetSequence() byte {
	return packet.Sequence
}

func (packet *Packet) String() string {
	if packet.IsComplete {
		return fmt.Sprintf("[Mariadb Packet: Length=%d Sequence=%d Body=%v]", packet.Length, packet.Sequence, packet.Body)
	} else {
		return fmt.Sprintf("[Mariadb Incomplete Packet: Length=%d Sequence=%d Body=%v]", packet.Length, packet.Sequence, packet.Body)
	}
}
