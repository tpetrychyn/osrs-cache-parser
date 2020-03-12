package compression

import (
	"golang.org/x/crypto/xtea"
	"hash"
	"io"
)

type NoneImpl struct {

}

func (n *NoneImpl) Decompress(reader io.Reader, compressedLength int32, crc hash.Hash32, xteaCipher *xtea.Cipher) ([]byte, error) {
	encryptedData := make([]byte, compressedLength)
	reader.Read(encryptedData)
	crc.Write(encryptedData)

	return encryptedData, nil
}
