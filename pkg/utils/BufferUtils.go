package utils

import (
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