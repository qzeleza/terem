package tui

import (
	"fmt"
	"time"

	"github.com/qzeleza/terem/internal/utils"
	"github.com/qzeleza/termos"
)

// SysInfoResult результат запуска DNS сервера
type SysInfoResult struct {
	Model       string        // Модель роутера
	Arch        string        // Архитектура
	MemoryUsage utils.RAMInfo // Использование памяти
	Uptime      time.Time     // Время работы
	Hostname    string        // Доменное имя
	IP          string        // IP-адрес
	Gateway     string        // Шлюз
	MAC         string        // MAC-адрес
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
				fmt.Sprintf("───────────────────"),
				fmt.Sprintf("Модель       : %s", result.Model),
				fmt.Sprintf("Архитектура  : %s", result.Arch),
				fmt.Sprintf("Память       : %d МБ / %d МБ (свободно: %d МБ)",
					result.MemoryUsage.Total-result.MemoryUsage.Free,
					result.MemoryUsage.Total,
					result.MemoryUsage.Free),
				fmt.Sprintf("Время работы : %s", utils.FormatUptime(result.Uptime)),
				fmt.Sprintf("Доменное имя : %s", result.Hostname),
				fmt.Sprintf("IP-адрес     : %s", result.IP),
				fmt.Sprintf("Шлюз         : %s", result.Gateway),
				fmt.Sprintf("MAC-адрес    : %s", result.MAC),
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
		MemoryUsage: utils.RAMInfo{Total: 0, Free: 0},
		Uptime:      time.Now(),
		Hostname:    "Неизвестно",
		IP:          "Неизвестно",
		Gateway:     "Неизвестно",
		MAC:         "Неизвестно",
	}

	// Получаем модель роутера
	if model, err := utils.GetRouterModel(); err == nil {
		result.Model = model
	}

	// Получаем архитектуру процессора
	if arch, err := utils.GetSystemArch(); err == nil {
		result.Arch = arch
	}

	// Получаем информацию о памяти
	if memInfo, err := utils.GetMemoryInfo(); err == nil {
		result.MemoryUsage = memInfo
	}

	// Получаем время работы системы
	if uptime, err := utils.GetSystemUptime(); err == nil {
		result.Uptime = uptime
	}

	// Получаем имя хоста
	if hostname, err := utils.GetHostname(); err == nil {
		result.Hostname = hostname
	}

	// Получаем информацию о сети
	if netInfo, err := utils.GetNetworkInfo(); err == nil {
		result.IP = netInfo.IP
		result.Gateway = netInfo.Gateway
		result.MAC = netInfo.MAC
	}
}
