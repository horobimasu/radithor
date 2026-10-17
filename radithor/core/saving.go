package core

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"radithor/store"

	"shared/config"

	"github.com/horobimasu/gotils"
)

type DiskFile struct {
	path string
	data []byte
}

func GetOrCreateFile(name string) (*DiskFile, error) {
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
			key := config.GetConfigAny[[]byte](store.DISK_FILE_DATA_ENCRYPTION_KEY)

			decrrypted, err := gotils.DecryptAesGcm(data, key)
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
		key := config.GetConfigAny[[]byte](store.DISK_FILE_DATA_ENCRYPTION_KEY)

		encrypted, err := gotils.EncryptAesGcm(data, key)
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
