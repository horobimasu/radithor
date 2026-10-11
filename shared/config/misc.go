package config

import (
	"debug/pe"
	"errors"

	"github.com/horobimasu/gotils"
)

var ErrNoConfig = errors.New("config does not exist")

var DecryptedConfigData = []byte{}

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

func RegisterConfig(rawConfig []byte) {
	DecryptedConfigData = rawConfig
}

// we are just gonna assume that it is set

func GetConfigAny[T any](id uint8) T {
	var thing T

	data, found := GetConfigRaw(DecryptedConfigData, id)
	if !found {
		return thing
	}

	value, _ := data.(T)
	return value
}

func GetConfigString(id uint8) string {
	data, found := GetConfigRaw(DecryptedConfigData, id)
	if !found {
		return ""
	}

	bytes, correct := data.([]byte)
	if !correct {
		return ""
	}

	return string(bytes)
}
