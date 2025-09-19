package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecuteCommand выполняет команду локально или удалённо
// В будущем можно расширить для поддержки SSH
func ExecuteCommand(cmd string) (string, error) {
	// Для локального выполнения используем sh
	command := exec.Command("sh", "-c", cmd)
	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out

	err := command.Run()
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения команды '%s': %v", cmd, err)
	}

	return strings.TrimSpace(out.String()), nil
}

// readFile читает содержимое файла
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла %s: %v", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}
