package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/qzeleza/terem/internal/i18n"
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
		return "", fmt.Errorf(i18n.T("utils.error.command"), cmd, err)
	}

	return strings.TrimSpace(out.String()), nil
}

// readFile читает содержимое файла
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf(i18n.T("utils.error.read_file"), path, err)
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
		return fmt.Sprintf(i18n.T("utils.duration.days_hours_minutes"), days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf(i18n.T("utils.duration.hours_minutes"), hours, minutes)
	} else {
		return fmt.Sprintf(i18n.T("utils.duration.minutes"), minutes)
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

// PadRight дополняет строку пробелами справа до указанной ширины
//
// Параметры:
//   - s: исходная строка
//   - width: желаемая ширина строки
//
// Возвращает:
//   - строка, дополненная пробелами до указанной ширины
//   - если исходная строка длиннее или равна width, возвращается исходная строка
func PadRight(s string, width int) string {
	// Если строка уже нужной длины или длиннее, возвращаем как есть
	if len(s) >= width {
		return s
	}

	// Создаем слайс байт нужной длины
	result := make([]byte, width)

	// Копируем исходную строку
	copy(result, s)

	// Заполняем оставшееся место пробелами
	for i := len(s); i < width; i++ {
		result[i] = ' '
	}

	return string(result)
}
