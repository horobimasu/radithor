package utils

import (
	"bytes"
	"encoding/binary"
)

func Deserialize[T any](rawData []byte) (T, error) {
	var result T

	err := binary.Read(
		bytes.NewReader(rawData),
		binary.LittleEndian,
		&result,
	)

	return result, err
}

func Serialize[T any](data T) ([]byte, error) {
	var result bytes.Buffer

	err := binary.Write(
		&result,
		binary.LittleEndian,
		data,
	)

	return result.Bytes(), err
}
