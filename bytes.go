package goutils

// writes zeroes to every byte in []byte
func ZeroBytes(b []byte) {
	for i, l := 0, len(b); i < l; i++ {
		b[i] = byte(0)
	}
}

// writes incrementing short as byte to []byte 0-255,0-255
func FillBytes(b []byte) {
	for i, l := 0, len(b); i < l; i++ {
		b[i] = byte(i)
	}
}

// return a count of the bytes that are diff
func DiffBytes(a, b []byte) int {
	ct, l := 0, 0
	if len(a) < len(b) {
		l = len(a)
	} else {
		l = len(b)
	}
	for i := 0; i < l; i++ {
		if a[i] != b[i] {
			ct++
		}
	}
	return ct
}
