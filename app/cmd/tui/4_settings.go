package tui

import (
	"github.com/qzeleza/termos"
)

// settingsList содержит список настроек приложения
var settingsList = []string{
	"Режим логирования",
	"Назад",
}

// SelectSettingsLoop отображает меню настроек приложения
func (ac *AppConfig) SelectSettingsLoop() {
	for {
		ac.SelectSettings()
		switch ac.Category {
		case settingsList[0]: // Режим логирования
			ac.SetDebugMode()
		case settingsList[1]: // Назад
			return
		default:
			ac.Log.Warn("Неверный выбор настройки")
			return
		}
	}
}

func (ac *AppConfig) SelectSettings() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор настроек приложения").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем задачу для выбора пункта меню
	menuTask := termos.NewSingleSelectTask("Выбор настроек", settingsList)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		err := ac.Log.Fatal("Ошибка при выборе настроек:", err)
		if err != nil {
			ac.Log.Error("Ошибка при выводе сообщения:", err)
		}
	}

	// ac.Settings = settingsList[menuTask.GetSelectedIndex()]
}

func (ac *AppConfig) SetDebugMode() {
	ac.Debug = !ac.Debug
	err := ac.Log.Info("Режим логирования:", ac.Debug)
	if err != nil {
		ac.Log.Error("Ошибка при выводе сообщения:", err)
	}
	err = ac.SetupLogger()
	if err != nil {
		ac.Log.Error("Ошибка при настройке логгера:", err)
	}
}
