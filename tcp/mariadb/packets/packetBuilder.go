package packets

import "fmt"

type PacketBuilder struct {
	DebugPrint bool
	buffer     []byte
}

func NewPacketBuilder() PacketBuilder {
	return PacketBuilder{
		DebugPrint: false,
		buffer:     make([]byte, 0),
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
		//		fmt.Printf("Buffersize < 4 (%d)\n\r", len(builder.buffer))
		return nil
	}

	//	length := int(binary.LittleEndian.Uint32([]byte{builder.buffer[0], builder.buffer[1], builder.buffer[2], 0x00}))
	length := int(builder.buffer[0]) + (int(builder.buffer[1]) << 8) + (int(builder.buffer[2]) << 16)
	packetSize := length + 4

	tmpSize := 40
	if len(builder.buffer) > tmpSize {
		builder.printf(" ReadPacket: %d (%d)\n\r  %v\n\r  ...\n\r  %v\n\r", packetSize, len(builder.buffer), builder.buffer[:tmpSize], builder.buffer[len(builder.buffer)-tmpSize:])
	} else {
		builder.printf(" ReadPacket: %d (%d)\n\r  %v\n\r", packetSize, len(builder.buffer), builder.buffer)
	}

	if len(builder.buffer) >= packetSize {
		pkg := PacketFromBytes(builder.buffer[:packetSize])
		builder.buffer = builder.buffer[packetSize:]

		//		_, file, no, ok := runtime.Caller(1)
		//		if ok {
		//			fmt.Printf("<<%d>> called from %s#%d\n\r", pkg.Sequence, file, no)
		//		}
		builder.printf("<<%d>> %s\r\n", pkg.Sequence, pkg.String())

		return pkg
	}
	return nil
}

func (builder *PacketBuilder) printf(msg string, a ...any) {
	if builder.DebugPrint {
		fmt.Printf(msg, a...)
	}
}
