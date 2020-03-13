package utils

import (
	"bytes"
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