package utils

func RSHashString(str string) int32 {
	var res int32
	for i:=0;i<len(str);i++ {
		res = (res << 5) - res + int32(CharToByteCp1252(rune(str[i])))
	}

	return res
}

func CharToByteCp1252(c rune) byte {
	if c > 0 && c < 128 || c >= 160 && c < 255 {
		return byte(c)
	} else if c == 8364 {
		return -128+256
	} else if c == 8218 {
		return -126+256
	} else if c == 402 {
		return -125+256
	} else if c == 8222 {
		return -124+256
	} else if c == 8230 {
		return -123+256
	} else if c == 8224 {
		return -122+256
	} else if c == 8225 {
		return -121+256
	} else if c == 710 {
		return -120+256
	} else if c == 8240 {
		return -119+256
	} else if c == 352 {
		return -118+256
	} else if c == 8249 {
		return -117+256
	} else if c == 338 {
		return -116+256
	} else if c == 381 {
		return -114+256
	} else if c == 8216 {
		return -111+256
	} else if c == 8217 {
		return -110+256
	} else if c == 8220 {
		return -109+256
	} else if c == 8221 {
		return -108+256
	} else if c == 8226 {
		return -107+256
	} else if c == 8211 {
		return -106+256
	} else if c == 8212 {
		return -105+256
	} else if c == 732 {
		return -104+256
	} else if c == 8482 {
		return -103+256
	} else if c == 353 {
		return -102+256
	} else if c == 8250 {
		return -101+256
	} else if c == 339 {
		return -100+256
	} else if c == 382 {
		return -98+256
	} else if c == 376 {
		return -97+256
	} else {
		return 63
	}
}
