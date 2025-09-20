package tui

import (
	"log"

	"github.com/qzeleza/terem/internal/i18n"
	"github.com/qzeleza/termos"
)

var mainMenuKeys = []string{
	ModeApps,
	ModeSettings,
	ModeExit,
}

// SelectMainMenu выбирает режим работы приложения
func (ac *AppConfig) SelectMainMenu() {

	// Устанавливаем язык для Термоса (TUI)
	termos.SetDefaultLanguage(i18n.Language())
	ac.AppTitle = i18n.T("app.title")

	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue(i18n.T("menu.main.queue.title")).
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Выводим информацию о системе
	ac.SysInfo(setupQueue)

	// Создаем список для выбора
	menuLabels := labelsFor(mainMenuKeys)

	// Создаем задачу для выбора пункта меню с запоминанием последней позиции
	menuTask := termos.NewSingleSelectTask(i18n.T("menu.main.task.title"), menuLabels).WithDefaultItem(ac.LastMainMenuIndex)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		log.Fatal(i18n.T("menu.main.error"), err)
	}
	// Сохраняем выбранный индекс и устанавливаем режим
	selected := menuTask.GetSelectedIndex()
	ac.LastMainMenuIndex = selected
	ac.Mode = mainMenuKeys[selected]
}

//
