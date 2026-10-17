package utils

import (
	"errors"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/horobimasu/gotils"
)

func GetRandomDir(
	rootPath string,
	maxDepth int,
	sampleSize int,
	blacklistedPaths []string,
) (string, error) {
	rootPath = gotils.ExpandPath(rootPath)

	var dirs = make([]string, 0, sampleSize)

	err := filepath.WalkDir(rootPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			if path == rootPath {
				return err
			}

			if entry != nil && entry.IsDir() {
				return filepath.SkipDir
			}

			return err
		}

		relative, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}

		relative = strings.ToLower(relative)
		path = strings.ToLower(path)

		for _, toSkip := range blacklistedPaths {
			if strings.HasPrefix(toSkip, "!") {
				exclude := strings.TrimPrefix(toSkip, "!")
				exclude = strings.ToLower(exclude)

				if exclude != "" && strings.Contains(relative, exclude) {
					return filepath.SkipDir
				}

				continue
			}

			toSkip = gotils.ExpandPath(toSkip)

			if strings.HasPrefix(toSkip, "*\\") {
				suffix := strings.TrimPrefix(toSkip, "*\\")
				if path == suffix || strings.HasSuffix(path, "\\"+suffix) {
					return filepath.SkipDir
				}

				continue
			}

			if relative == toSkip || strings.HasPrefix(relative, toSkip+"\\") {
				return filepath.SkipDir
			}

			if path == toSkip || strings.HasPrefix(path, toSkip+"\\") {
				return filepath.SkipDir
			}
		}

		depth := len(strings.Split(relative, "\\"))
		if depth > maxDepth {
			return filepath.SkipDir
		}

		if depth != maxDepth {
			return nil
		}

		testFile, err := os.CreateTemp(path, "z.q*")
		if err != nil {
			return filepath.SkipDir
		}

		err = testFile.Close()
		if err != nil {
			os.Remove(testFile.Name())
			return filepath.SkipDir
		}

		err = os.Remove(testFile.Name())
		if err != nil {
			return filepath.SkipDir
		}

		dirs = append(dirs, path)
		if len(dirs) >= sampleSize {
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if len(dirs) == 0 {
		return "", errors.New("failed to find any directories")
	}

	return dirs[rand.Intn(len(dirs))], nil
}
