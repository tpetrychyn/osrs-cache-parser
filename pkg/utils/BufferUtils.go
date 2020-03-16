package utils

import (
	"bytes"
	"encoding/binary"
	"io"
)

func ReadString(r io.Reader) string {
	result := ""
	for {
		var b = make([]byte, 1)
		_, err := r.Read(b)
		if err != nil {
			return ""
		}
		if b[0] == 0 {
			break
		}
		result += string(b)
	}

	return result
}

func Read24BitInt(reader *bytes.Reader) int32 {
	by := make([]byte, 3)
	reader.Read(by)
	return int32(by[0])<<16 + int32(by[1])<<8 + int32(by[2])
}

func ReadShortSmart(buf *bytes.Reader) (int, error) {
	peek, _ := buf.ReadByte()
	buf.UnreadByte()
	peek = peek & 0xFF
	if peek < 128 {
		b, err := buf.ReadByte()
		return int(b) - 64, err
	} else {
		var short uint16
		err := binary.Read(buf, binary.BigEndian, &short)
		return int(short) - 49152, err
	}
}