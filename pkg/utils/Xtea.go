package utils

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/xtea"
	"io/ioutil"
)

func XteaDecrypt(cipher *xtea.Cipher, src []byte) []byte {
	numBlocks := len(src) / xtea.BlockSize
	res := make([]byte, 0, len(src))

	// pad to an even block size of 8
	if len(src) % xtea.BlockSize != 0 {
		pad := xtea.BlockSize - len(src) % xtea.BlockSize
		src = append(src, make([]byte, pad)...)
	}

	// iterate 8 bytes at a time decrypting them - one block at a time
	for i:=0;i<numBlocks;i++ {
		piece := src[i*xtea.BlockSize:(i+1)*xtea.BlockSize]
		dec := make([]byte, len(piece))
		cipher.Decrypt(dec, piece)
		res = append(res, dec...)
	}

	return res
}

func XteaKeyFromIntArray(keys []int32) (*xtea.Cipher, error) {
	xteaKey := make([]byte, 16)
	for i := 0; i < len(keys); i++ {
		j := i << 2
		xteaKey[j] = byte(keys[i] >> 24)
		xteaKey[j+1] = byte(keys[i] >> 16)
		xteaKey[j+2] = byte(keys[i] >> 8)
		xteaKey[j+3] = byte(keys[i])
	}

	xteaCipher, err := xtea.NewCipher(xteaKey)
	if err != nil {
		return nil, err
	}
	return xteaCipher, nil
}

func LoadXteas() (map[uint16][]int32, error) {
	var xteaDefs = make(map[uint16][]int32)
	file, err := ioutil.ReadFile("../../cache/xteas.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open xteas.json %w", err)
	}

	type XteaDef struct {
		Region uint16
		Keys []int32
	}

	var xteas []*XteaDef
	err = json.Unmarshal(file, &xteas)
	if err != nil {
		return nil, fmt.Errorf("failed to parse xteas.json %w", err)
	}

	for _, v := range xteas {
		xteaDefs[v.Region] = v.Keys
	}

	return xteaDefs, nil
}