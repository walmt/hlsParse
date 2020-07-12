package util

func BytesIsZero(buf []byte) bool {
	for i := 0; i < len(buf); i++ {
		if buf[i] != 0 {
			return false
		}
	}
	return true
}

func BytesIsOxFF(buf []byte) bool {
	for i := 0; i < len(buf); i++ {
		if buf[i] != 0xFF {
			return false
		}
	}
	return true
}
