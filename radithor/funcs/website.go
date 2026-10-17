package funcs

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/horobimasu/gotils"
)

const _SITE_BLOCK_IP = "127.0.0.22"

func BlockWebsite(website string) (string, error) { // the cleaned site it blocked, err
	hostName, err := cleanWebsite(website)
	if err != nil {
		return "", err
	}

	hostPath := gotils.ExpandPath("%SYSTEMROOT%\\System32\\drivers\\etc\\hosts")

	data, err := os.ReadFile(hostPath)
	if err != nil {
		return "", err
	}

	for line := range strings.Lines(string(data)) {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 && strings.EqualFold(fields[1], hostName) {
			return "", errors.New("site already blocked")
		}
	}

	file, err := os.OpenFile(hostPath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return "", err
	}

	defer file.Close()

	file.WriteString(fmt.Sprintf("\n%s %s", _SITE_BLOCK_IP, hostName))

	return hostName, gotils.RunCommand("ipconfig", "/flushdns")
}

func UnblockWebsite(website string) (string, error) { // the cleaned site it unblocked, err
	hostName, err := cleanWebsite(website)
	if err != nil {
		return "", err
	}

	hostPath := gotils.ExpandPath("%SYSTEMROOT%\\System32\\drivers\\etc\\hosts")

	data, err := os.ReadFile(hostPath)
	if err != nil {
		return "", err
	}

	var newData strings.Builder

	for line := range strings.Lines(string(data)) {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			newData.WriteString(line)
			continue
		}

		fields := strings.Fields(trimmed)
		if len(fields) >= 2 && strings.EqualFold(fields[1], hostName) {
			continue
		}

		newData.WriteString(line)
	}

	err = os.WriteFile(hostPath, []byte(newData.String()), 0666)
	if err != nil {
		return "", err
	}

	return hostName, gotils.RunCommand("ipconfig", "/flushdns")
}

func ClearAllBlockedSites() ([]string, error) { // all the sites that where unblocked, err
	hostPath := gotils.ExpandPath("%SYSTEMROOT%\\System32\\drivers\\etc\\hosts")

	data, err := os.ReadFile(hostPath)
	if err != nil {
		return nil, err
	}

	var newData strings.Builder
	var unblocked []string

	for line := range strings.Lines(string(data)) {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == _SITE_BLOCK_IP {
			unblocked = append(unblocked, fields[1])
			continue
		}

		newData.WriteString(line)
	}

	err = os.WriteFile(hostPath, []byte(newData.String()), 0666)
	if err != nil {
		return nil, err
	}

	return unblocked, gotils.RunCommand("ipconfig", "/flushdns")
}

func cleanWebsite(website string) (string, error) {
	if !strings.Contains(website, "://") {
		website = "//" + website
	}

	parsed, err := url.Parse(website)
	if err != nil {
		return "", err
	}

	return parsed.Hostname(), nil
}
