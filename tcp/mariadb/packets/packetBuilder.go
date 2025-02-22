package packets

import (
	"encoding/binary"
)

type PacketBuilder struct {
	buffer []byte
}

func NewPacketBuilder() *PacketBuilder {
	return &PacketBuilder{
		buffer: make([]byte, 0),
	}
}

func (builder *PacketBuilder) AddBytes(data []byte) {
	builder.buffer = append(builder.buffer, data...)
}

func (builder *PacketBuilder) GetBuffer() []byte {
	return builder.buffer
}

func (builder *PacketBuilder) BuildPacket() *Packet {
	if len(builder.buffer) < 4 {
		return nil
	}

	length := int(binary.LittleEndian.Uint32([]byte{builder.buffer[0], builder.buffer[1], builder.buffer[2], 0x00}))
	packetSize := length + 4

	if len(builder.buffer) >= packetSize {
		pkg := PacketFromBytes(builder.buffer[:packetSize])
		builder.buffer = builder.buffer[packetSize:]
		return pkg
	}
	return nil
}
