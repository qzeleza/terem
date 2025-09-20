package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// networkInfo структура для хранения сетевой информации
type networkInfo struct {
	IP      string
	Gateway string
	MAC     string
}

type RAMInfo struct {
	Total int // Общее количество памяти
	Free  int // Свободная память
}

// GetRouterModel получает модель роутера
func GetRouterModel() (string, error) {
	// Сначала пробуем получить из /proc/cpuinfo
	if content, err := ReadFile("/proc/cpuinfo"); err == nil {
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "machine") || strings.HasPrefix(line, "Model") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
		}
	}

	// Пробуем получить через ubus (специфично для OpenWrt)
	if output, err := ExecuteCommand("ubus call system board 2>/dev/null | grep -E '\"model\"' | cut -d'\"' -f4"); err == nil && output != "" {
		return output, nil
	}

	// Пробуем через файл board.json
	if content, err := ReadFile("/etc/board.json"); err == nil {
		// Простой парсинг для извлечения модели
		re := regexp.MustCompile(`"model"\s*:\s*"([^"]+)"`)
		if matches := re.FindStringSubmatch(content); len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("не удалось определить модель роутера")
}

// GetSystemArch получает архитектуру процессора
func GetSystemArch() (string, error) {
	// Используем uname -m для получения архитектуры
	arch, err := ExecuteCommand("uname -m")
	if err != nil {
		return "", fmt.Errorf("ошибка получения архитектуры: %v", err)
	}
	return arch, nil
}

// GetMemoryInfo получает информацию о памяти
func GetMemoryInfo() (RAMInfo, error) {
	var memInfo RAMInfo

	// Читаем /proc/meminfo
	content, err := ReadFile("/proc/meminfo")
	if err != nil {
		return memInfo, fmt.Errorf("ошибка чтения информации о памяти: %v", err)
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if val, err := strconv.Atoi(fields[1]); err == nil {
					memInfo.Total = val / 1024 // Конвертируем из KB в MB
				}
			}
		} else if strings.HasPrefix(line, "MemFree:") || strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if val, err := strconv.Atoi(fields[1]); err == nil {
					// Используем MemAvailable если есть, иначе MemFree
					if strings.HasPrefix(line, "MemAvailable:") || memInfo.Free == 0 {
						memInfo.Free = val / 1024 // Конвертируем из KB в MB
					}
				}
			}
		}
	}

	if memInfo.Total == 0 {
		return memInfo, fmt.Errorf("не удалось получить информацию о памяти")
	}

	return memInfo, nil
}

// GetSystemUptime получает время работы системы
func GetSystemUptime() (time.Time, error) {
	// Читаем /proc/uptime
	content, err := ReadFile("/proc/uptime")
	if err != nil {
		return time.Time{}, fmt.Errorf("ошибка чтения uptime: %v", err)
	}

	// Первое число - время работы в секундах
	fields := strings.Fields(content)
	if len(fields) == 0 {
		return time.Time{}, fmt.Errorf("неверный формат uptime")
	}

	// Парсим секунды с плавающей точкой
	uptimeStr := strings.Split(fields[0], ".")[0]
	uptimeSeconds, err := strconv.ParseInt(uptimeStr, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("ошибка парсинга uptime: %v", err)
	}

	// Вычисляем время запуска системы
	bootTime := time.Now().Add(-time.Duration(uptimeSeconds) * time.Second)
	return bootTime, nil
}

// GetHostname получает имя хоста
func GetHostname() (string, error) {
	// Пробуем прочитать из /proc/sys/kernel/hostname
	if hostname, err := ReadFile("/proc/sys/kernel/hostname"); err == nil && hostname != "" {
		return hostname, nil
	}

	// Альтернативный способ через команду hostname
	if hostname, err := ExecuteCommand("hostname"); err == nil && hostname != "" {
		return hostname, nil
	}

	// Используем стандартный Go метод как запасной вариант
	if hostname, err := os.Hostname(); err == nil {
		return hostname, nil
	}

	return "", fmt.Errorf("не удалось получить имя хоста")
}

// GetNetworkInfo получает сетевую информацию
func GetNetworkInfo() (*networkInfo, error) {
	info := &networkInfo{
		IP:      "Неизвестно",
		Gateway: "Неизвестно",
		MAC:     "Неизвестно",
	}

	// Получаем основной сетевой интерфейс и шлюз
	defaultIface := ""
	if output, err := ExecuteCommand("ip route show default 2>/dev/null | head -1"); err == nil {
		// Парсим строку вида: default via 192.168.1.1 dev br-lan ...
		fields := strings.Fields(output)
		for i, field := range fields {
			if field == "via" && i+1 < len(fields) {
				info.Gateway = fields[i+1]
			}
			if field == "dev" && i+1 < len(fields) {
				defaultIface = fields[i+1]
			}
		}
	}

	// Если не нашли интерфейс через route, пробуем найти первый активный
	if defaultIface == "" {
		// Получаем список интерфейсов
		if output, err := ExecuteCommand("ip link show up 2>/dev/null | grep -E '^[0-9]+:' | grep -v 'lo:' | head -1"); err == nil {
			// Парсим строку вида: 2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP>...
			fields := strings.Fields(output)
			if len(fields) >= 2 {
				defaultIface = strings.TrimSuffix(fields[1], ":")
			}
		}
	}

	// Получаем IP-адрес и MAC-адрес интерфейса
	if defaultIface != "" {
		// Получаем IP-адрес
		if output, err := ExecuteCommand(fmt.Sprintf("ip addr show %s 2>/dev/null | grep 'inet ' | head -1", defaultIface)); err == nil {
			// Парсим строку вида: inet 192.168.1.100/24 brd ...
			fields := strings.Fields(output)
			if len(fields) >= 2 {
				ip := strings.Split(fields[1], "/")[0]
				if ip != "" {
					info.IP = ip
				}
			}
		}

		// Получаем MAC-адрес
		if output, err := ExecuteCommand(fmt.Sprintf("ip link show %s 2>/dev/null | grep 'link/ether' | head -1", defaultIface)); err == nil {
			// Парсим строку вида: link/ether 00:11:22:33:44:55 brd ...
			fields := strings.Fields(output)
			if len(fields) >= 2 {
				info.MAC = fields[1]
			}
		}
	}

	// Альтернативный способ через ifconfig (для старых систем)
	if info.IP == "Неизвестно" || info.MAC == "Неизвестно" {
		if output, err := ExecuteCommand("ifconfig 2>/dev/null | grep -A1 'inet addr' | head -2"); err == nil {
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				// Парсим IP
				if strings.Contains(line, "inet addr:") {
					re := regexp.MustCompile(`inet addr:([0-9.]+)`)
					if matches := re.FindStringSubmatch(line); len(matches) > 1 {
						info.IP = matches[1]
					}
				}
				// Парсим MAC
				if strings.Contains(line, "HWaddr") {
					re := regexp.MustCompile(`HWaddr ([0-9A-Fa-f:]+)`)
					if matches := re.FindStringSubmatch(line); len(matches) > 1 {
						info.MAC = matches[1]
					}
				}
			}
		}
	}

	return info, nil
}
