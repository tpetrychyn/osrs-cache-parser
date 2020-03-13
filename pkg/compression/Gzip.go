package compression

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
	"golang.org/x/crypto/xtea"
	"hash"
	"io"
)

type Gzip struct {}

func (g *Gzip) Decompress(reader io.Reader, compressedLength int32, crc hash.Hash32, xteaCipher *xtea.Cipher) ([]byte, error) {
	encryptedData := make([]byte, compressedLength+4)
	reader.Read(encryptedData)
	crc.Write(encryptedData)

	if xteaCipher != nil {
		encryptedData = utils.XteaDecrypt(xteaCipher, encryptedData)
	}

	stream := bytes.NewReader(encryptedData)

	var decompressedLength int32
	binary.Read(stream, binary.BigEndian, &decompressedLength)

	comp := make([]byte, compressedLength)
	stream.ReadAt(comp, 4)

	gz, err := gzip.NewReader(bytes.NewReader(comp))
	if err != nil {
		panic(err)
	}

	var dBuffer bytes.Buffer
	if _, err := io.Copy(&dBuffer, gz); err != nil {
		panic(err)
	}

	if len(dBuffer.Bytes()) != int(decompressedLength) {
		return nil, fmt.Errorf("bytes read from bzip %d did not match decompressedLength %d", len(dBuffer.Bytes()), decompressedLength)
	}

	return dBuffer.Bytes(), nil
}
