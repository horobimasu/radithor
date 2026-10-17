package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

func CorruptFile(filePath string, passes uint) error {
	if passes < 1 {
		return errors.New("passes cannot be less than 1")
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return errors.New("the provided file is not regular")
	}

	fileSize := info.Size()
	if fileSize == 0 {
		return nil
	}

	for i := range passes {
		_, err := file.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to seek on pass %d", i)
		}

		written, err := io.CopyN(file, rand.Reader, fileSize)
		if err != nil {
			return fmt.Errorf("failed to write on pass %d", i)
		}

		if written != fileSize {
			return fmt.Errorf("failed to finish writing on pass %d", i)
		}

		err = file.Sync()
		if err != nil {
			return err
		}
	}

	return nil
}
