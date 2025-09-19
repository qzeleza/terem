// Package router предоставляет функционал для работы с роутерами
package utils

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Router представляет информацию о роутере
type Router struct {
	Name     string
	Address  string
	Arch     string
	Platform string
	SSHPort  string
}

// Utility содержит информацию о системной утилите
type Utility struct {
	Name    string
	Package string
}

// Список утилит, необходимых для работы с роутером
var (
	// RequiredUtilities содержит список обязательных утилит и их пакетов в OpenWrt
	RequiredUtilities = []Utility{
		{Name: "curl", Package: "curl"},
		{Name: "wget", Package: "wget"},
		{Name: "nc", Package: "ncat-ssl"},
		{Name: "ipset", Package: "ipset"},
		{Name: "iptables", Package: "iptables"},
	}

	// defaultSSHPort порт SSH по умолчанию
	defaultSSHPort = "22"
)

// RunCommand выполняет команду на удаленном роутере через SSH
func (r Router) RunCommand(command string) (string, error) {
	if r.SSHPort == "" {
		r.SSHPort = defaultSSHPort
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"ssh",
		"-p", r.SSHPort,
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		fmt.Sprintf("root@%s", r.Address),
		command,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ошибка выполнения команды: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// CheckRequiredUtilities проверяет наличие обязательных утилит на роутере
// и при необходимости устанавливает их
func (r Router) CheckRequiredUtilities() []string {
	var results []string

	for _, util := range RequiredUtilities {
		// Проверяем наличие утилиты
		_, err := r.RunCommand(fmt.Sprintf("which %s", util.Name))
		if err == nil {
			results = append(results, fmt.Sprintf("%s - уже установлена", util.Name))
			continue
		}

		// Пытаемся установить утилиту
		_, err = r.RunCommand(fmt.Sprintf("opkg update && opkg install %s", util.Package))
		if err != nil {
			results = append(results, fmt.Sprintf("%s - ошибка установки: %v", util.Name, err))
			continue
		}

		// Проверяем установку
		if _, err := r.RunCommand(fmt.Sprintf("which %s", util.Name)); err != nil {
			results = append(results, fmt.Sprintf("%s - не удалось установить", util.Name))
			continue
		}

		results = append(results, fmt.Sprintf("%s - успешно установлена", util.Name))
	}

	return results
}

// GetSystemInfo возвращает базовую информацию о системе роутера
func (r Router) GetSystemInfo() (map[string]string, error) {
	info := make(map[string]string)

	// Получаем информацию о системе
	output, err := r.RunCommand("uname -a")
	if err != nil {
		return nil, fmt.Errorf("не удалось получить информацию о системе: %v", err)
	}
	info["system_info"] = output

	// Получаем информацию о памяти
	output, err = r.RunCommand("free -m")
	if err != nil {
		return nil, fmt.Errorf("не удалось получить информацию о памяти: %v", err)
	}
	info["memory"] = output

	// Получаем информацию о диске
	output, err = r.RunCommand("df -h")
	if err != nil {
		return nil, fmt.Errorf("не удалось получить информацию о диске: %v", err)
	}
	info["disk"] = output

	return info, nil
}
