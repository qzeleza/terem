package tui

import (
	"github.com/qzeleza/terem/internal/i18n"
	"github.com/qzeleza/termos"
)

// networkList содержит список сетевых приложений в фиксированном порядке
var networkList = []string{
	NetworkOptionOpenSSH,
	NetworkOptionProxy,
	NetworkOptionDNS,
	NetworkOptionAdGuard,
	NetworkOptionBack,
}

// networkKeys соответствующие ключи для networkList
var networkKeys = []string{
	"openssh",
	"proxy",
	"dns",
	"adguard",
	"quit",
}

// NetworkCategoryLoop запускает цикл для выбора сетевых приложений
func (ac *AppConfig) NetworkCategoryLoop() {
	ac.ContextualLoop(func() bool {
		ac.SelectNetworkCategory()

		// Проверяем контекст после выбора
		if ac.IsContextCancelled() {
			return false
		}

		switch ac.Category {
		case NetworkOptionOpenSSH:
			ac.SelectOpenSSHApp()
			return true
		case NetworkOptionProxy:
			ac.SelectProxyApp()
			return true
		case NetworkOptionDNS:
			ac.SelectDNSApp()
			return true
		case NetworkOptionAdGuard:
			ac.SelectAdGuardApp()
			return true
		case NetworkOptionBack:
			return false
		default:
			ac.Log.Warn(i18n.T("network.warn.invalid"))
			return false
		}
	}, i18n.T("loop.network"))
}

// SelectNetworkCategory отображает меню для выбора сетевых приложений
func (ac *AppConfig) SelectNetworkCategory() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue(i18n.T("network.queue.title")).
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	labels := labelsFor(networkList)

	// Создаем задачу для выбора пункта меню с запоминанием последней позиции
	menuTask := termos.NewSingleSelectTask(i18n.T("network.task.title"), labels).WithDefaultItem(ac.LastNetworkIndex)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		ac.Log.Fatal(i18n.T("network.error"), err)
	}

	// Если пользователь отменил выбор, возвращаемся к предыдущему меню
	if menuTask.HasError() {
		ac.Category = NetworkOptionBack
		return
	}

	// Сохраняем выбранный индекс и устанавливаем категорию
	selected := menuTask.GetSelectedIndex()
	ac.LastNetworkIndex = selected
	ac.Category = networkList[selected]
}

// SelectOpenSSHApp отображает меню для выбора OpenSSH-сервера
func (ac *AppConfig) SelectOpenSSHApp() {
	ac.Log.Info(i18n.T("network.log.openssh"))
}

// SelectProxyApp отображает меню для выбора прокси сервера
func (ac *AppConfig) SelectProxyApp() {
	ac.Log.Info(i18n.T("network.log.proxy"))
}

// SelectDNSApp отображает меню для выбора DNS-сервера
func (ac *AppConfig) SelectDNSApp() {
	ac.Log.Info(i18n.T("network.log.dns"))
}

// SelectAdGuardApp отображает меню для выбора AdGuard Home сервера
func (ac *AppConfig) SelectAdGuardApp() {
	ac.Log.Info(i18n.T("network.log.adguard"))
}
