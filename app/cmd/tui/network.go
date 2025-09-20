package tui

import (
	"github.com/qzeleza/termos"
)

// networkList содержит список сетевых приложений в фиксированном порядке
var networkList = []string{
	"OpenSSH-сервер",
	"Прокси сервер 3proxy",
	"DNSmasq-сервер",
	"AdGuard Home сервер",
	"Назад",
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
		case networkList[0]: // OpenSSH-сервер
			ac.SelectOpenSSHApp()
			// После выполнения действия показываем меню снова
			return true
		case networkList[1]: // Прокси сервер 3proxy
			ac.SelectProxyApp()
			// После выполнения действия показываем меню снова
			return true
		case networkList[2]: // DNSmasq-сервер
			ac.SelectDNSApp()
			// После выполнения действия показываем меню снова
			return true
		case networkList[3]: // AdGuard Home сервер
			ac.SelectAdGuardApp()
			// После выполнения действия показываем меню снова
			return true
		case networkList[4]: // Назад
			return false
		default:
			ac.Log.Warn("Неверный выбор категории")
			return false
		}
	}, "цикла сетевых приложений")
}

// SelectNetworkCategory отображает меню для выбора сетевых приложений
func (ac *AppConfig) SelectNetworkCategory() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор сетевых приложений").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем задачу для выбора пункта меню с запоминанием последней позиции
	menuTask := termos.NewSingleSelectTask("Выберите сетевое приложение", networkList).WithDefaultItem(ac.LastNetworkIndex)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		ac.Log.Fatal("Ошибка при выборе сетевого приложения:", err)
	}

	// Сохраняем выбранный индекс и устанавливаем категорию
	ac.LastNetworkIndex = menuTask.GetSelectedIndex()
	ac.Category = networkList[menuTask.GetSelectedIndex()]
}

// SelectOpenSSHApp отображает меню для выбора OpenSSH-сервера
func (ac *AppConfig) SelectOpenSSHApp() {
	ac.Log.Info("Выбран OpenSSH-сервер")
}

// SelectProxyApp отображает меню для выбора прокси сервера
func (ac *AppConfig) SelectProxyApp() {
	ac.Log.Info("Выбран прокси сервер 3proxy")
}

// SelectDNSApp отображает меню для выбора DNS-сервера
func (ac *AppConfig) SelectDNSApp() {
	ac.Log.Info("Выбран DNSmasq-сервер")
}

// SelectAdGuardApp отображает меню для выбора AdGuard Home сервера
func (ac *AppConfig) SelectAdGuardApp() {
	ac.Log.Info("Выбран AdGuard Home сервер")
}
