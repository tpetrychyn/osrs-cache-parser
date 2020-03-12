package utils

func Djb2(s string) int32 {

	var hash int32 = 0
	for _, c := range s {
		hash = c + ((hash << 5) - hash)
	}

	return hash
}


