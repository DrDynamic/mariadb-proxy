package readers

import (
	"encoding/binary"
)

type LeByteReader struct {
	data []byte
}

func NewLeByteReader(data []byte) *LeByteReader {
	return &LeByteReader{
		data: data,
	}
}

func (reader *LeByteReader) PopBytes(count int) []byte {
	data := reader.data[:count]
	reader.data = reader.data[count:]

	return data
}

func (reader *LeByteReader) PopNullTerminatedBytes() []byte {
	terminatorIndex := 0
	for ; reader.data[terminatorIndex] != 0x00; terminatorIndex++ {
	}

	data := reader.PopBytes(terminatorIndex)
	reader.PopBytes(1) // pop nullterminator
	return data
}

func (reader *LeByteReader) PopLengthEncodedBytes() []byte {
	length := reader.PopLengthEncodedInt()
	return reader.PopBytes(int(length))
}

func (reader *LeByteReader) PopUInt8() uint8 {
	return reader.PopBytes(1)[0]
}

func (reader *LeByteReader) PopUInt16() uint16 {
	return binary.LittleEndian.Uint16(reader.PopBytes(2))
}

func (reader *LeByteReader) PopUInt24() uint32 {
	data := reader.PopBytes(3)
	return binary.LittleEndian.Uint32([]byte{data[0], data[1], data[2], 0x00})
}

func (reader *LeByteReader) PopUInt32() uint32 {
	return binary.LittleEndian.Uint32(reader.PopBytes(4))
}

func (reader *LeByteReader) PopUInt64() uint64 {
	return binary.LittleEndian.Uint64(reader.PopBytes(8))
}

func (reader *LeByteReader) PopLengthEncodedInt() uint64 {
	first := reader.PopUInt8()
	if first < 0xFB {
		return uint64(first)
	} else if first == 0xFB {
		return 0
	} else if first == 0xFC {
		return uint64(reader.PopUInt16())
	} else if first == 0xFD {
		return uint64(reader.PopUInt24())
	} else if first == 0xFE {
		return reader.PopUInt64()
	}
	return 0
}

func (reader *LeByteReader) PopString(count int) string {
	return string(reader.PopBytes(count))
}

func (reader *LeByteReader) PopNullTerminatedString() string {
	terminatorIndex := 0
	for ; reader.data[terminatorIndex] != 0x00; terminatorIndex++ {
	}

	data := string(reader.PopBytes(terminatorIndex))
	reader.PopBytes(1) // pop nullterminator
	return data
}

func (reader *LeByteReader) PopLengthEncodedString() string {
	length := reader.PopLengthEncodedInt()
	return reader.PopString(int(length))
}

func (reader *LeByteReader) PopEofEncodedString() string {
	return reader.PopString(reader.Length())
}

func (reader *LeByteReader) Length() int {
	return len(reader.data)
}
