package tui

import (
	"fmt"
	"time"

	"github.com/qzeleza/terem/internal/i18n"
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
	task := termos.NewFuncTask(i18n.T("sysinfo.task.title"),
		func() error {
			ac.Log.Info(i18n.T("sysinfo.log.fetch"))
			// Используем кешированные данные
			_ = ac.GetSysInfo()
			return nil
		},
		termos.WithSummaryFunction(func() []string {
			// Получаем кешированные данные для отображения
			info := ac.GetSysInfo()
			divider := "────────────────────────────"
			maxLength := 15
			return []string{
				divider,
				fmt.Sprintf("%s: %s", utils.PadRight(i18n.T("sysinfo.summary.model"), maxLength), info.Model),
				fmt.Sprintf("%s: %s", utils.PadRight(i18n.T("sysinfo.summary.arch"), maxLength), info.Arch),
				fmt.Sprintf("%s: %d/%d/%d Mb", utils.PadRight(i18n.T("sysinfo.summary.memory"), maxLength),
					info.MemoryUsage.Total-info.MemoryUsage.Free,
					info.MemoryUsage.Total,
					info.MemoryUsage.Free),
				fmt.Sprintf("%s: %s", utils.PadRight(i18n.T("sysinfo.summary.uptime"), maxLength), utils.FormatUptime(info.Uptime)),
				fmt.Sprintf("%s: %s", utils.PadRight(i18n.T("sysinfo.summary.hostname"), maxLength), info.Hostname),
				fmt.Sprintf("%s: %s", utils.PadRight(i18n.T("sysinfo.summary.ip"), maxLength), info.IP),
				fmt.Sprintf("%s: %s", utils.PadRight(i18n.T("sysinfo.summary.gateway"), maxLength), info.Gateway),
				fmt.Sprintf("%s: %s", utils.PadRight(i18n.T("sysinfo.summary.mac"), maxLength), info.MAC),
				// divider,
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
		Model:       i18n.T("sysinfo.default"),
		Arch:        i18n.T("sysinfo.default"),
		MemoryUsage: utils.RAMInfo{Total: 0, Free: 0},
		Uptime:      time.Now(),
		Hostname:    i18n.T("sysinfo.default"),
		IP:          i18n.T("sysinfo.default"),
		Gateway:     i18n.T("sysinfo.default"),
		MAC:         i18n.T("sysinfo.default"),
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
