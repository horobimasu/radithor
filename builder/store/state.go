package store

import (
	"os"
	"path/filepath"
)

var (
	SharedSecret []byte

	BuilderFilePath = ""
	BuilderDirPath  = ""
)

func InitState() {
	BuilderFilePath, _ = os.Executable()

	BuilderDirPath = filepath.Dir(BuilderFilePath)
}
