package tui

import (
	"github.com/qzeleza/terem/internal/i18n"
	"github.com/qzeleza/termos"
)

// settingsList содержит список настроек приложения
var settingsList = []string{
	SettingsOptionLogging,
	SettingsOptionBack,
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
		case SettingsOptionLogging:
			ac.SetDebugMode()
			return true
		case SettingsOptionBack:
			return false
		default:
			ac.Log.Warn(i18n.T("settings.warn.invalid"))
			return false
		}
	}, i18n.T("loop.settings"))
}

func (ac *AppConfig) SelectSettings() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue(i18n.T("settings.queue.title")).
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	labels := labelsFor(settingsList)

	// Создаем задачу для выбора пункта меню с запоминанием последней позиции
	menuTask := termos.NewSingleSelectTask(i18n.T("settings.task.title"), labels).WithDefaultItem(ac.LastSettingsIndex)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		ac.Log.Fatal(i18n.T("settings.error"), err)
	}

	// Сохраняем выбранный индекс и устанавливаем категорию
	selected := menuTask.GetSelectedIndex()
	ac.LastSettingsIndex = selected
	ac.Category = settingsList[selected]
}

func (ac *AppConfig) SetDebugMode() {
	ac.Debug = !ac.Debug
	ac.Conf.DebugMode = ac.Debug
	ac.Log.Info(i18n.T("settings.log.toggle"), ac.Debug)
	if err := ac.SetupLogger(); err != nil {
		ac.Log.Fatal(i18n.T("cli.debug.error"), err)
	}
}
