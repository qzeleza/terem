package zlog

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// formatArgs форматирует аргументы для логирования.
func formatArgs(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	if format, ok := args[0].(string); ok && len(args) > 1 {
		return fmt.Sprintf(format, args[1:]...)
	}
	return fmt.Sprint(args...)
}

// getTotalMemory возвращает общий объём физической памяти в байтах (для Linux/embedded).
// Читает /proc/meminfo, ищет "MemTotal:" (стандарт для роутеров на Linux).
// Если ошибка — возвращает 0 (тогда используются стандартные настройки).
func getTotalMemory() uint64 {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0 // Ошибка — fallback на defaults.
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			var memKB uint64
			_, err := fmt.Sscanf(line, "MemTotal: %d kB", &memKB)
			if err != nil {
				return 0
			}
			return memKB * 1024 // Конвертируем kB в байты.
		}
	}
	return 0 // Не найдено — fallback.
}
