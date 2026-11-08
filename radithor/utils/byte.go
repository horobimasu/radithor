package utils

import (
	"encoding/binary"

	"golang.org/x/sys/windows"
)

func StringToUTF16LeBytes(str string) []byte {
	encoded := windows.StringToUTF16(str)

	buffer := make([]byte, len(encoded)*2)

	for i, char := range encoded {
		binary.LittleEndian.PutUint16(buffer[i*2:], char)
	}

	return buffer
}

func AlignBytes(bytes []byte, alignment int) []byte {
	for len(bytes)%alignment != 0 {
		bytes = append(bytes, 0)
	}

	return bytes
}
