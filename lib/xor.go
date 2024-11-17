package lib

func XorBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		panic("длины массивов должны совпадать")
	}
	result := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}
	return result
}
