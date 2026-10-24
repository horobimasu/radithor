package store

import (
	"os"
	"path/filepath"
)

var (
	SharedSecret []byte

	BuilderFilePath = ""
	BuilderDirPath  = ""

	BuildId int64
)

func InitState() {
	BuilderFilePath, _ = os.Executable()

	BuilderDirPath = filepath.Dir(BuilderFilePath)
}
