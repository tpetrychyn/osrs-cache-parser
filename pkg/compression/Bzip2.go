package compression

import (
	"bytes"
	"compress/bzip2"
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/xtea"
	"hash"
	"io"
)

var BzipHeader = []byte{66, 90, 104, 49}

type Bzip2 struct {}

func (b *Bzip2) Decompress(reader io.Reader, compressedLength int32, crc hash.Hash32, xteaCipher *xtea.Cipher) ([]byte, error) {
	encryptedData := make([]byte, compressedLength+4)
	_, _ = reader.Read(encryptedData)
	crc.Write(encryptedData)

	stream := bytes.NewReader(encryptedData)
	var decompressedLength int32
	err := binary.Read(stream, binary.BigEndian, &decompressedLength)
	if err != nil {
		return nil, err
	}

	comp := make([]byte, compressedLength)
	_, err = stream.ReadAt(comp, 4)
	if err != nil {
		return nil, err
	}

	comp = append(BzipHeader, comp...)
	bz := bzip2.NewReader(bytes.NewReader(comp))

	var dBuffer bytes.Buffer
	if _, err := io.Copy(&dBuffer, bz); err != nil {
		return nil, err
	}

	if len(dBuffer.Bytes()) != int(decompressedLength) {
		return nil, fmt.Errorf("bytes read from bzip %d did not match decompressedLength %d", len(dBuffer.Bytes()), decompressedLength)
	}

	return dBuffer.Bytes(), nil
}
