package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

// FormatUptime форматирует время работы системы в читаемый вид
func FormatUptime(bootTime time.Time) string {
	// Вычисляем продолжительность работы
	uptime := time.Since(bootTime)

	// Извлекаем дни, часы, минуты
	days := int(uptime.Hours() / 24)
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60

	// Форматируем строку
	if days > 0 {
		return fmt.Sprintf("%d дн. %d ч. %d мин.", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d ч. %d мин.", hours, minutes)
	} else {
		return fmt.Sprintf("%d мин.", minutes)
	}
}

// GetEnv получает значение переменной окружения или возвращает значение по умолчанию.
func GetEnv(key, defaultValue string) string {
	// Получаем значение переменной окружения
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetEnvBool получает значение переменной окружения или возвращает значение по умолчанию.
func GetEnvBool(key string, defaultValue bool) bool {
	// Получаем значение переменной окружения
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	// Преобразуем значение в булево
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}
