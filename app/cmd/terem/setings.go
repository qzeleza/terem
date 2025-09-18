package terem

import (
	"log"

	"github.com/qzeleza/termos"
)

var settingsList = []string{
	"Режим логирования",
	"Вернуться",
}

// SelectSettings отображает меню настроек приложения
func (ac *AppConfig) SelectSettings() {
	for {
		ac.SelectSettingsFromList()
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

func (ac *AppConfig) SelectSettingsFromList() {
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
		log.Fatal("Ошибка при выборе настроек:", err)
	}

	// ac.Settings = settingsList[menuTask.GetSelectedIndex()]
}

func (ac *AppConfig) SetDebugMode() {
	ac.Debug = !ac.Debug
	ac.Log.Info("Режим логирования:", ac.Debug)
	ac.SetupLogger()
}
