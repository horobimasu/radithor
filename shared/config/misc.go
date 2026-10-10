package config

import (
	"debug/pe"
	"errors"

	"github.com/horobimasu/gotils"
)

var ErrNoConfig = errors.New("config does not exist")

func GetAndDecryptConfig(execFilePath string, sectionName string, encryptionKey []byte) ([]byte, error) {
	file, err := pe.Open(execFilePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	section := file.Section(sectionName)
	if section == nil {
		return nil, ErrNoConfig
	}

	data, err := section.Data()
	if err != nil {
		return nil, err
	}

	if int(section.VirtualSize) < len(data) {
		data = data[:section.VirtualSize] // the thing that adds sections will pad the section
	}

	return gotils.DecryptAesGcm(data, encryptionKey)
}
