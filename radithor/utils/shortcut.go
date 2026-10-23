package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/horobimasu/gotils"
)

type Shortcut struct {
	ShortcutPath string
	Target       string
	IconPath     string
	Args         string
	StartIn      string

	// snca
	Description string
	Key         string
	Style       int
}

func CreateShortcut(shortcutData Shortcut) error {
	var command strings.Builder

	command.WriteString("$k = New-Object -ComObject WScript.Shell\n")
	fmt.Fprintf(&command, "$i = $k.CreateShortcut(%q)\n", gotils.ExpandPath(shortcutData.ShortcutPath))
	fmt.Fprintf(&command, "$i.TargetPath = %q\n", shortcutData.Target)
	fmt.Fprintf(&command, "$i.IconLocation = %q\n", shortcutData.IconPath)
	fmt.Fprintf(&command, "$i.Arguments = %q\n", shortcutData.Args)
	fmt.Fprintf(&command, "$i.WorkingDirectory = %q\n", shortcutData.StartIn)
	fmt.Fprintf(&command, "$i.Description = %q\n", shortcutData.Description)
	fmt.Fprintf(&command, "$i.HotKey = %q\n", shortcutData.Key)
	fmt.Fprintf(&command, "$i.WindowStyle = %d\n", shortcutData.Style)
	command.WriteString("$i.Save()")

	return gotils.RunCommand("powershell", "-nop", "-ep", "bypass", "-c", command.String())
}

func GetShortcutData(shortcutPath string) (Shortcut, error) {
	shortcutPath = gotils.ExpandPath(shortcutPath)

	result := Shortcut{
		ShortcutPath: shortcutPath,
	}

	var command strings.Builder

	command.WriteString("$u = New-Object -ComObject WScript.Shell\n")
	fmt.Fprintf(&command, "$c = $u.CreateShortcut(%q)\n", shortcutPath)
	fmt.Fprintf(&command, "Write-Host \"2,\" $c.TargetPath\n")
	fmt.Fprintf(&command, "Write-Host \"3,\" $c.IconLocation\n")
	fmt.Fprintf(&command, "Write-Host \"4,\" $c.Arguments\n")
	fmt.Fprintf(&command, "Write-Host \"5,\" $c.WorkingDirector\n")
	fmt.Fprintf(&command, "Write-Host \"6,\" $c.Description\n")
	fmt.Fprintf(&command, "Write-Host \"7,\" $c.HotKey\n")
	fmt.Fprintf(&command, "Write-Host \"9,\" $c.WindowStyle\n")

	out, err := gotils.GetCommandOutput("powershell", "-nop", "-ep", "bypass", "-c", command.String())
	if err != nil {
		return result, err
	}

	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "2,"):
			result.Target = strings.TrimSpace(strings.TrimPrefix(line, "2,"))
		case strings.HasPrefix(line, "3,"):
			result.IconPath = strings.TrimSpace(strings.TrimPrefix(line, "3,"))
		case strings.HasPrefix(line, "4,"):
			result.Args = strings.TrimSpace(strings.TrimPrefix(line, "4,"))
		case strings.HasPrefix(line, "5,"):
			result.StartIn = strings.TrimSpace(strings.TrimPrefix(line, "5,"))
		case strings.HasPrefix(line, "6,"):
			result.Description = strings.TrimSpace(strings.TrimPrefix(line, "6,"))
		case strings.HasPrefix(line, "7,"):
			result.Key = strings.TrimSpace(strings.TrimPrefix(line, "7,"))
		case strings.HasPrefix(line, "9,"):
			result.Style, _ = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "9,")))
		}
	}

	return result, nil
}

func UpdateShortcutData(shortcutPath string, newData Shortcut) error {
	shortcutPath = gotils.ExpandPath(shortcutPath)

	data, err := GetShortcutData(shortcutPath)
	if err != nil {
		return err
	}

	shortcutData := Shortcut{
		ShortcutPath: shortcutPath,
		Target:       data.Target,
		IconPath:     data.IconPath,
		Args:         data.Args,
		StartIn:      data.StartIn,
		Description:  data.Description,
		Key:          data.Key,
		Style:        data.Style,
	}

	if newData.Target != "" {
		shortcutData.Target = newData.Target
	}

	if newData.IconPath != "" {
		shortcutData.IconPath = newData.IconPath
	}

	if newData.Args != "" {
		shortcutData.Args = newData.Args
	}

	if newData.StartIn != "" {
		shortcutData.StartIn = newData.StartIn
	}

	if newData.Description != "" {
		shortcutData.Description = newData.Description
	}

	if newData.Key != "" {
		shortcutData.Key = newData.Key
	}

	if newData.Style != 0 {
		shortcutData.Style = newData.Style
	}

	return CreateShortcut(shortcutData)
}
