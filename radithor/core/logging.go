package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"radithor/store"

	"github.com/horobimasu/gotils"
)

const _LOGS_FILE_NAME = "logs"

var (
	logDiskFile      *DiskFile
	logProxyFilePath string

	userprofilePath string
)

func LogInfo(text string, args ...any) {
	log(fmt.Sprintf("[INFO] "+text, args...))
}

func LogWarning(text string, args ...any) {
	log(fmt.Sprintf("[WARNING] "+text, args...))
}

func LogError(text string, args ...any) {
	log(fmt.Sprintf("[ERROR] "+text, args...))
}

func log(text string) {
	if store.DISALLOW_LOGGING {
		return
	}

	pid := os.Getpid()

	now := time.Now()
	var time any = now.Unix()

	if store.DEBUG_MODE_ENABLED {
		time = strings.ToLower(now.Format("02/01/2006 3:04 PM"))
	}

	diskFile, proxyFilePath, err := GetLogsFilePath()
	if err != nil {
		ExitFatalOnDebugMode(fmt.Sprintf("failed to get or create debug logs file - %v", err))
		return
	}

	line := fmt.Sprintf("[%d] [%s] %s\n", pid, time, sanitizeLogText(text))

	err = diskFile.AppendDataAndSave([]byte(line))
	if err != nil {
		ExitFatalOnDebugMode(fmt.Sprintf("failed to write to log file - %v", err))
		logDiskFile = nil // maybe resetting it will fix whatever the error is
	}

	if proxyFilePath != "" {
		err := os.MkdirAll(filepath.Dir(proxyFilePath), 0755)
		if err != nil {
			ExitFatalOnDebugMode(fmt.Sprintf("failed to create dir for proxy logs file - %v", err))
			return
		}

		file, err := os.OpenFile(proxyFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			ExitFatalOnDebugMode(fmt.Sprintf("failed to open proxy logs file - %v", err))
			return
		}

		defer file.Close()

		fmt.Fprint(file, line)
	}
}

func GetLogsFilePath() (*DiskFile, string, error) { // disk file, proxy file path, error
	if logDiskFile == nil {
		file, err := GetOrCreateDiskFile(_LOGS_FILE_NAME)
		if err != nil {
			return nil, "", err
		}

		logDiskFile = file
	}

	if logProxyFilePath == "" && store.DEBUG_MODE_ENABLED && store.DEBUG_MODE_LOG_FILE_PROXY_TO_PATH != "" {
		logProxyFilePath = gotils.ExpandPath(store.DEBUG_MODE_LOG_FILE_PROXY_TO_PATH)
	}

	return logDiskFile, logProxyFilePath, nil
}

func sanitizeLogText(text string) string {
	text = strings.ReplaceAll(text, "\\\\", "\\")

	if store.SHORTEN_USERPROFILE_PATHS_IN_LOGS {
		if userprofilePath == "" {
			userprofilePath = os.Getenv("USERPROFILE")
		}

		if userprofilePath != "" {
			text = strings.ReplaceAll(text, userprofilePath, "~")
		}
	}

	return text
}
