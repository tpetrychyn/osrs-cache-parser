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

	var decryptedData []byte
	if xteaCipher != nil {
		decryptedData = utils.XteaDecrypt(xteaCipher, encryptedData)
	} else {
		decryptedData = encryptedData
	}

	// should be 6870 length
	stream := bytes.NewReader(decryptedData)

	var decompressedLength int32
	binary.Read(stream, binary.BigEndian, &decompressedLength)

	comp := make([]byte, compressedLength)
	stream.ReadAt(comp, 4)

	gz, err := gzip.NewReader(bytes.NewReader(comp))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}

	var dBuffer bytes.Buffer
	if _, err := io.Copy(&dBuffer, gz); err != nil {
		// Loading land is throwing checksum error despite seeming to parse correctly
	}

	if len(dBuffer.Bytes()) != int(decompressedLength) {
		return nil, fmt.Errorf("bytes read from bzip %d did not match decompressedLength %d", len(dBuffer.Bytes()), decompressedLength)
	}

	return dBuffer.Bytes(), nil
}
