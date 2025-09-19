package tui

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/qzeleza/terem/internal/utils"
	"github.com/qzeleza/termos"
)

type RAMInfo struct {
	Total int // Общее количество памяти
	Free  int // Свободная память
}

// SysInfoResult результат запуска DNS сервера
type SysInfoResult struct {
	Model       string    // Модель роутера
	Arch        string    // Архитектура
	MemoryUsage RAMInfo   // Использование памяти
	Uptime      time.Time // Время работы
	Hostname    string    // Доменное имя
	IP          string    // IP-адрес
	Route       string    // Шлюз
	MAC         string    // MAC-адрес
}

// networkInfo структура для хранения сетевой информации
type networkInfo struct {
	ip      string
	gateway string
	mac     string
}

// SysInfo выводит информацию о системе
func (ac *AppConfig) SysInfo(queue *termos.Queue) {
	// Создаем переменную для хранения результата
	var result SysInfoResult

	// Выводим информацию о системе
	task := termos.NewFuncTask("Информация о системе",
		func() error {
			ac.Log.Info("Получение информации о системе")
			getSysInfo(&result)
			return nil
		},
		termos.WithSummaryFunction(func() []string {
			return []string{
				fmt.Sprintf("Модель: %s", result.Model),
				fmt.Sprintf("Архитектура: %s", result.Arch),
				fmt.Sprintf("Память: %d МБ / %d МБ (свободно: %d МБ)",
					result.MemoryUsage.Total-result.MemoryUsage.Free,
					result.MemoryUsage.Total,
					result.MemoryUsage.Free),
				fmt.Sprintf("Время работы: %s", formatUptime(result.Uptime)),
				fmt.Sprintf("Доменное имя: %s", result.Hostname),
				fmt.Sprintf("IP-адрес: %s", result.IP),
				fmt.Sprintf("Шлюз: %s", result.Route),
				fmt.Sprintf("MAC-адрес: %s", result.MAC),
			}
		}),
		termos.WithStopOnError(true),
	)

	queue.AddTasks(task)
}

// getSysInfo получает информацию о системе роутера
func getSysInfo(result *SysInfoResult) {
	// Инициализируем структуру с значениями по умолчанию
	*result = SysInfoResult{
		Model:       "Неизвестно",
		Arch:        "Неизвестно",
		MemoryUsage: RAMInfo{Total: 0, Free: 0},
		Uptime:      time.Now(),
		Hostname:    "Неизвестно",
		IP:          "Неизвестно",
		Route:       "Неизвестно",
		MAC:         "Неизвестно",
	}

	// Получаем модель роутера
	if model, err := getRouterModel(); err == nil {
		result.Model = model
	}

	// Получаем архитектуру процессора
	if arch, err := getSystemArch(); err == nil {
		result.Arch = arch
	}

	// Получаем информацию о памяти
	if memInfo, err := getMemoryInfo(); err == nil {
		result.MemoryUsage = memInfo
	}

	// Получаем время работы системы
	if uptime, err := getSystemUptime(); err == nil {
		result.Uptime = uptime
	}

	// Получаем имя хоста
	if hostname, err := getHostname(); err == nil {
		result.Hostname = hostname
	}

	// Получаем информацию о сети
	if netInfo, err := getNetworkInfo(); err == nil {
		result.IP = netInfo.ip
		result.Route = netInfo.gateway
		result.MAC = netInfo.mac
	}
}

// getRouterModel получает модель роутера
func getRouterModel() (string, error) {
	// Сначала пробуем получить из /proc/cpuinfo
	if content, err := utils.ReadFile("/proc/cpuinfo"); err == nil {
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
	if output, err := utils.ExecuteCommand("ubus call system board 2>/dev/null | grep -E '\"model\"' | cut -d'\"' -f4"); err == nil && output != "" {
		return output, nil
	}

	// Пробуем через файл board.json
	if content, err := utils.ReadFile("/etc/board.json"); err == nil {
		// Простой парсинг для извлечения модели
		re := regexp.MustCompile(`"model"\s*:\s*"([^"]+)"`)
		if matches := re.FindStringSubmatch(content); len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("не удалось определить модель роутера")
}

// getSystemArch получает архитектуру процессора
func getSystemArch() (string, error) {
	// Используем uname -m для получения архитектуры
	arch, err := utils.ExecuteCommand("uname -m")
	if err != nil {
		return "", fmt.Errorf("ошибка получения архитектуры: %v", err)
	}
	return arch, nil
}

// getMemoryInfo получает информацию о памяти
func getMemoryInfo() (RAMInfo, error) {
	var memInfo RAMInfo

	// Читаем /proc/meminfo
	content, err := utils.ReadFile("/proc/meminfo")
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

// getSystemUptime получает время работы системы
func getSystemUptime() (time.Time, error) {
	// Читаем /proc/uptime
	content, err := utils.ReadFile("/proc/uptime")
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

// getHostname получает имя хоста
func getHostname() (string, error) {
	// Пробуем прочитать из /proc/sys/kernel/hostname
	if hostname, err := utils.ReadFile("/proc/sys/kernel/hostname"); err == nil && hostname != "" {
		return hostname, nil
	}

	// Альтернативный способ через команду hostname
	if hostname, err := utils.ExecuteCommand("hostname"); err == nil && hostname != "" {
		return hostname, nil
	}

	// Используем стандартный Go метод как запасной вариант
	if hostname, err := os.Hostname(); err == nil {
		return hostname, nil
	}

	return "", fmt.Errorf("не удалось получить имя хоста")
}

// getNetworkInfo получает сетевую информацию
func getNetworkInfo() (*networkInfo, error) {
	info := &networkInfo{
		ip:      "Неизвестно",
		gateway: "Неизвестно",
		mac:     "Неизвестно",
	}

	// Получаем основной сетевой интерфейс и шлюз
	defaultIface := ""
	if output, err := utils.ExecuteCommand("ip route show default 2>/dev/null | head -1"); err == nil {
		// Парсим строку вида: default via 192.168.1.1 dev br-lan ...
		fields := strings.Fields(output)
		for i, field := range fields {
			if field == "via" && i+1 < len(fields) {
				info.gateway = fields[i+1]
			}
			if field == "dev" && i+1 < len(fields) {
				defaultIface = fields[i+1]
			}
		}
	}

	// Если не нашли интерфейс через route, пробуем найти первый активный
	if defaultIface == "" {
		// Получаем список интерфейсов
		if output, err := utils.ExecuteCommand("ip link show up 2>/dev/null | grep -E '^[0-9]+:' | grep -v 'lo:' | head -1"); err == nil {
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
		if output, err := utils.ExecuteCommand(fmt.Sprintf("ip addr show %s 2>/dev/null | grep 'inet ' | head -1", defaultIface)); err == nil {
			// Парсим строку вида: inet 192.168.1.100/24 brd ...
			fields := strings.Fields(output)
			if len(fields) >= 2 {
				ip := strings.Split(fields[1], "/")[0]
				if ip != "" {
					info.ip = ip
				}
			}
		}

		// Получаем MAC-адрес
		if output, err := utils.ExecuteCommand(fmt.Sprintf("ip link show %s 2>/dev/null | grep 'link/ether' | head -1", defaultIface)); err == nil {
			// Парсим строку вида: link/ether 00:11:22:33:44:55 brd ...
			fields := strings.Fields(output)
			if len(fields) >= 2 {
				info.mac = fields[1]
			}
		}
	}

	// Альтернативный способ через ifconfig (для старых систем)
	if info.ip == "Неизвестно" || info.mac == "Неизвестно" {
		if output, err := utils.ExecuteCommand("ifconfig 2>/dev/null | grep -A1 'inet addr' | head -2"); err == nil {
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				// Парсим IP
				if strings.Contains(line, "inet addr:") {
					re := regexp.MustCompile(`inet addr:([0-9.]+)`)
					if matches := re.FindStringSubmatch(line); len(matches) > 1 {
						info.ip = matches[1]
					}
				}
				// Парсим MAC
				if strings.Contains(line, "HWaddr") {
					re := regexp.MustCompile(`HWaddr ([0-9A-Fa-f:]+)`)
					if matches := re.FindStringSubmatch(line); len(matches) > 1 {
						info.mac = matches[1]
					}
				}
			}
		}
	}

	return info, nil
}

// formatUptime форматирует время работы системы в читаемый вид
func formatUptime(bootTime time.Time) string {
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
