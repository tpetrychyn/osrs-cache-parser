package compression

import (
	"golang.org/x/crypto/xtea"
	"hash"
	"io"
)

const (
	None int8 = iota
	BZIP2
	GZIP
)

type Compression interface {
	Decompress(reader io.Reader, compressedLength int32, crc hash.Hash32, xteaCipher *xtea.Cipher) ([]byte, error)
}

func GetCompressionStrategy(typ int8) Compression {
	switch typ {
	case None:
		return &NoneImpl{}
	case BZIP2:
		return &Bzip2{}
	case GZIP:
		return &Gzip{}
	}

	return &NoneImpl{}
}