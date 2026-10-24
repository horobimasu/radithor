package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	buildsDir := filepath.Join("..", ".builds")

	err := os.MkdirAll(buildsDir, 0775)
	if err != nil {
		return "", err
	}

	items, err := os.ReadDir(buildsDir)
	if err != nil {
		return "", err
	}

	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		if !strings.HasPrefix(item.Name(), "latest-") {
			continue
		}

		newName := strings.TrimPrefix(item.Name(), "latest-")

		oldPath := filepath.Join(buildsDir, item.Name())
		newPath := filepath.Join(buildsDir, newName)

		err := os.Rename(oldPath, newPath)
		if err != nil {
			return "", err
		}

		break // there will always only be 1
	}

	newDirPath := filepath.Join(buildsDir, fmt.Sprintf("latest-%d", store.BuildId))

	err = os.MkdirAll(newDirPath, 0755)
	if err != nil {
		return "", err
	}

	writeTo := filepath.Join(newDirPath, name)

	err = os.WriteFile(writeTo, data, 0666)
	if err != nil {
		return "", err
	}

	return writeTo, nil
}
