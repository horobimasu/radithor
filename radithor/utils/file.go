package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"radithor/core"
	"radithor/store"

	"github.com/horobimasu/gotils"
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

func CopyFile(srcPath string, dstPath string) error {
	err := gotils.CopyFile(srcPath, dstPath)
	if err != nil {
		return err
	}

	if !store.DEBUG_MODE_ENABLED || srcPath != store.ExecPath {
		return nil
	}

	hostName := filepath.Base(store.ExecPath)
	pattern := regexp.MustCompile("^" + regexp.QuoteMeta(hostName) + "_ads_[^_]+_[0-9a-f]{12}$")

	items, err := os.ReadDir(store.ExecDir)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.IsDir() {
			continue
		}

		if !pattern.MatchString(item.Name()) {
			continue
		}

		theSrcPath := filepath.Join(store.ExecDir, item.Name())
		theDstPath := filepath.Join(filepath.Dir(dstPath), item.Name())

		err := gotils.CopyFile(theSrcPath, theDstPath)
		if err != nil {
			return err
		}

		// if we are coping the rad file we can assume that
		// we are moving it and all the ads' need to be deleted anyway
		err = os.Remove(filepath.Join(store.ExecDir, item.Name()))
		if err != nil {
			core.LogWarning("ffailled to delete %s - %v", item.Name(), err) // todo, remove
		}
	}

	return nil
}
