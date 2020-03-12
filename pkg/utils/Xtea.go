package utils

import "golang.org/x/crypto/xtea"

func XteaDecrypt(cipher *xtea.Cipher, src []byte) []byte {
	// pad to an even block size of 8
	if len(src) % xtea.BlockSize != 0 {
		pad := xtea.BlockSize - len(src) % xtea.BlockSize
		src = append(src, make([]byte, pad)...)
	}

	// allocate an empty slice with a capacity of src length
	res := make([]byte, 0, len(src))

	// iterate 8 bytes at a time decrypting them - one block at a time
	for i:=0;i<len(src);i+=xtea.BlockSize {
		piece := src[i:i+xtea.BlockSize]
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