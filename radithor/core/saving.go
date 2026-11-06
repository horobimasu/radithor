package core

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"radithor/store"

	"github.com/horobimasu/gotils"
	"github.com/horobimasu/windowsx"
	"golang.org/x/sys/windows"
)

var (
	ErrEncryptionKeyNotSet = errors.New("encryption key not yet set")
	ErrFileBlocked         = errors.New("this disk file has been marked as blocked")
)

type DiskFile struct {
	path string
	data []byte
}

func GetAllDiskFilesNames() ([]string, error) {
	var names []string

	if store.DEBUG_MODE_ENABLED {
		hostName := filepath.Base(store.ExecPath)

		items, err := os.ReadDir(store.ExecDir)
		if err != nil {
			return nil, err
		}

		pattern := regexp.MustCompile("^" + regexp.QuoteMeta(hostName) + "_ads_[^_]+_[0-9a-f]{12}$")

		for _, item := range items {
			if item.IsDir() {
				continue
			}

			if !pattern.MatchString(item.Name()) {
				continue
			}

			_, after, _ := strings.Cut(item.Name(), "_ads_")
			names = append(names, strings.Split(after, "_")[0])
		}
	} else {
		var name string

		handle, err := windowsx.FindFirstStream(
			windows.StringToUTF16Ptr(store.ExecPath),
			windowsx.FindStreamInfoStandard,
			&name,
			0,
		)

		if err != nil {
			return nil, err
		}

		defer windows.FindClose(windows.Handle(handle))

		for {
			names = append(names, name)

			err = windowsx.FindNextStream(handle, &name)
			if err == windows.ERROR_HANDLE_EOF {
				break
			}

			if err != nil {
				return nil, err
			}
		}
	}

	return names, nil
}

func GetOrCreateDiskFile(name string) (*DiskFile, error) {
	// only the logger should get this error
	if !store.DiskFileDataEncryptionKeySet {
		return nil, ErrEncryptionKeyNotSet
	}

	fullHash := sha256.Sum256([]byte(name))
	nameHash := hex.EncodeToString(fullHash[:])[:12]

	var adsPath string

	if store.DEBUG_MODE_ENABLED {
		hostName := filepath.Base(store.ExecPath)
		adsName := fmt.Sprintf("%s_ads_%s_%s", hostName, name, nameHash)

		adsPath = filepath.Join(store.ExecDir, adsName)
	} else {
		adsPath = store.ExecPath + ":" + nameHash
	}

	data, err := os.ReadFile(adsPath)
	if err == nil {
		if !store.DEBUG_MODE_ENABLED {
			decrrypted, err := gotils.DecryptAesGcm(data, store.DiskFileDataEncryptionKey[:])
			if err != nil {
				return nil, err
			}

			data = decrrypted
		}

		return &DiskFile{
			path: adsPath,
			data: data,
		}, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	err = os.MkdirAll(filepath.Dir(adsPath), 0755)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(adsPath, nil, 0666)
	if err != nil {
		return nil, err
	}

	return &DiskFile{
		path: adsPath,
		data: []byte{},
	}, nil
}

func (file *DiskFile) GetData() []byte {
	return file.data
}

func (file *DiskFile) GetDataAsLines() []string {
	return strings.Split(string(file.data), "\n")
}

func (file *DiskFile) WriteDataAndSave(data []byte) error {
	file.data = data

	if !store.DEBUG_MODE_ENABLED {
		encrypted, err := gotils.EncryptAesGcm(data, store.DiskFileDataEncryptionKey[:])
		if err != nil {
			return err
		}

		data = encrypted
	}

	return os.WriteFile(file.path, data, 0666)
}

func (file *DiskFile) AppendDataAndSave(data []byte) error {
	return file.WriteDataAndSave(append(file.data, data...))
}
