package utils

import (
	"os"
	"path/filepath"

	"radithor-builder/store"
)

func GetSameDirFilePath(name string) string {
	path := filepath.Join(store.BuilderDirPath, name)

	_, err := os.Stat(path)
	if err != nil {
		return ""
	}

	return path
}

func WriteToBuildsDir(name string, data []byte) (string, error) { // where it wrote to, err
	err := os.MkdirAll("./../.builds", 0755)
	if err != nil {
		return "", err
	}

	writeTo := filepath.Join("./../.builds", name)

	err = os.WriteFile(writeTo, data, 0666)
	if err != nil {
		return "", err
	}

	return writeTo, nil
}
