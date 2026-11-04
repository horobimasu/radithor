package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"radithor/store"

	"github.com/horobimasu/gotils"
)

type DownloadProgress struct {
	Percent int
	Err     error
	Path    string
	Name    string
}

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

		newName := filepath.Base(dstPath) + strings.TrimPrefix(item.Name(), hostName)

		theSrcPath := filepath.Join(store.ExecDir, item.Name())
		theDstPath := filepath.Join(filepath.Dir(dstPath), newName)

		err := gotils.CopyFile(theSrcPath, theDstPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func DownloadFile(url string, dstDir string) (string, error) {
	return internalDownloadFile(url, dstDir, nil)
}

func DownloadFileWithProgress(url string, dstDir string) <-chan DownloadProgress {
	progress := make(chan DownloadProgress, 1)

	go func() {
		defer close(progress)

		update := func(downloadProgress DownloadProgress) {
			select {
			case progress <- downloadProgress:
			default:
				select {
				case <-progress:
				default:
				}

				progress <- downloadProgress
			}
		}

		downloadedTo, err := internalDownloadFile(url, dstDir, update)
		if err != nil {
			update(DownloadProgress{
				Err: err,
			})

			return
		}

		update(DownloadProgress{
			Percent: 100,
			Path:    downloadedTo,
		})
	}()

	return progress
}

func internalDownloadFile(
	url string,
	dstDir string,
	updateFunc func(DownloadProgress),
) (string, error) {
	dstDir = gotils.ExpandPath(dstDir)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK {
		return "", fmt.Errorf("bad response: %s", res.Status)
	}

	testFile, err := os.CreateTemp(dstDir, "p.u*")
	if err != nil {
		return "", err
	}

	testFilePath := testFile.Name()

	var wrote int64
	lastPercent := -2

	buffer := make([]byte, 32*1024)

	for {
		at, err := res.Body.Read(buffer)
		if at > 0 {
			written, err := testFile.Write(buffer[:at])
			wrote += int64(written)

			if err != nil {
				testFile.Close()
				return "", err
			}

			if written != at {
				testFile.Close()
				return "", io.ErrShortWrite
			}

			if updateFunc != nil {
				percent := -1

				if res.ContentLength > 0 {
					percent = int(wrote * 100 / res.ContentLength)
					percent = min(percent, 99)
				}

				if percent != lastPercent {
					updateFunc(DownloadProgress{
						Percent: percent,
					})

					lastPercent = percent
				}
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			testFile.Close()
			return "", err
		}
	}

	err = testFile.Close()
	if err != nil {
		return "", err
	}

	name, _, _ := strings.Cut(url, "?")
	extension := filepath.Ext(name)

	dstFilePath := filepath.Join(
		dstDir,
		fmt.Sprintf("%d%s", time.Now().UnixNano(), extension),
	)

	err = os.Rename(testFilePath, dstFilePath)
	if err != nil {
		return "", err
	}

	return dstFilePath, nil
}
