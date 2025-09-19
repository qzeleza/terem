package tui

import (
	"github.com/qzeleza/termos"
)

// networkList содержит список сетевых приложений
var networkList = map[string]string{
	"proxy":   "Прокси сервер 3proxy",
	"dns":     "DNSmasq-сервер",
	"adguard": "AdGuard Home сервер",
	"quit":    "Назад",
}

// NetworkCategoryLoop запускает цикл для выбора сетевых приложений
func (ac *AppConfig) NetworkCategoryLoop() {
	for {
		ac.SelectNetworkCategory()
		switch ac.Category {
		case networkList["proxy"]: // Прокси сервер 3proxy
			ac.SelectProxyApp()
		case networkList["dns"]: // DNSmasq-сервер
			ac.SelectDNSApp()
		case networkList["adguard"]: // AdGuard Home сервер
			ac.SelectAdGuardApp()
		case networkList["exit"]: // Выход
			return
		default:
			err := ac.Log.Warn("Неверный выбор категории")
			if err != nil {
				ac.Log.Error("Ошибка при выводе сообщения:", err)
			}
			return
		}
	}
}

// SelectNetworkCategory отображает меню для выбора сетевых приложений
func (ac *AppConfig) SelectNetworkCategory() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор сетевых приложений").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем список для выбора
	list := []string{}
	for _, v := range networkList {
		list = append(list, v)
	}

	// Создаем задачу для выбора пункта меню
	menuTask := termos.NewSingleSelectTask("Выберите сетевое приложение", list)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		err := ac.Log.Fatal("Ошибка при выборе сетевого приложения:", err)
		if err != nil {
			ac.Log.Error("Ошибка при выводе сообщения:", err)
		}
	}

	ac.Category = list[menuTask.GetSelectedIndex()]
}

// SelectProxyApp отображает меню для выбора прокси сервера
func (ac *AppConfig) SelectProxyApp() {
	err := ac.Log.Info("Выбран прокси сервер 3proxy")
	if err != nil {
		ac.Log.Error("Ошибка при выводе сообщения:", err)
	}
}

// SelectDNSApp отображает меню для выбора DNS-сервера
func (ac *AppConfig) SelectDNSApp() {
	err := ac.Log.Info("Выбран DNSmasq-сервер")
	if err != nil {
		ac.Log.Error("Ошибка при выводе сообщения:", err)
	}
}

// SelectAdGuardApp отображает меню для выбора AdGuard Home сервера
func (ac *AppConfig) SelectAdGuardApp() {
	err := ac.Log.Info("Выбран AdGuard Home сервер")
	if err != nil {
		ac.Log.Error("Ошибка при выводе сообщения:", err)
	}
}
