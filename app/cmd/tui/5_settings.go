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
	ac.ContextualLoop(func() bool {
		ac.SelectSettings()

		// Проверяем контекст после выбора
		if ac.IsContextCancelled() {
			return false
		}

		switch ac.Category {
		case settingsList[0]: // Режим логирования
			ac.SetDebugMode()
			// После выполнения действия показываем меню снова
			return true
		case settingsList[1]: // Назад
			return false
		default:
			ac.Log.Warn("Неверный выбор настройки")
			return false
		}
	}, "настроек")
}

func (ac *AppConfig) SelectSettings() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор настроек приложения").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем задачу для выбора пункта меню с запоминанием последней позиции
	menuTask := termos.NewSingleSelectTask("Выбор настроек", settingsList).WithDefaultItem(ac.LastSettingsIndex)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		ac.Log.Fatal("Ошибка при выборе настроек:", err)
	}

	// Сохраняем выбранный индекс и устанавливаем категорию
	ac.LastSettingsIndex = menuTask.GetSelectedIndex()
	ac.Category = settingsList[menuTask.GetSelectedIndex()]
}

func (ac *AppConfig) SetDebugMode() {
	ac.Debug = !ac.Debug
	ac.Conf.DebugMode = ac.Debug
	ac.Log.Info("Режим логирования:", ac.Debug)
	if err := ac.SetupLogger(); err != nil {
		ac.Log.Fatal("Ошибка при настройке логгера:", err)
	}
}
