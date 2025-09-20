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

	// Выводим информацию о системе
	task := termos.NewFuncTask("Информация о системе",
		func() error {
			ac.Log.Info("Получение информации о системе")
			// Используем кешированные данные
			_ = ac.GetSysInfo()
			return nil
		},
		termos.WithSummaryFunction(func() []string {
			// Получаем кешированные данные для отображения
			info := ac.GetSysInfo()
			return []string{
				fmt.Sprintf("%s", "───────────────────"),
				fmt.Sprintf("Модель       : %s", info.Model),
				fmt.Sprintf("Архитектура  : %s", info.Arch),
				fmt.Sprintf("Память       : %d МБ / %d МБ (свободно: %d МБ)",
					info.MemoryUsage.Total-info.MemoryUsage.Free,
					info.MemoryUsage.Total,
					info.MemoryUsage.Free),
				fmt.Sprintf("Время работы : %s", utils.FormatUptime(info.Uptime)),
				fmt.Sprintf("Доменное имя : %s", info.Hostname),
				fmt.Sprintf("IP-адрес     : %s", info.IP),
				fmt.Sprintf("Шлюз         : %s", info.Gateway),
				fmt.Sprintf("MAC-адрес    : %s", info.MAC),
			}
		}),
		termos.WithStopOnError(true),
	)

	queue.AddTasks(task)
}

// getSysInfo получает информацию о системе роутера
func (ac *AppConfig) getSysInfo(result *SysInfoResult) {
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
