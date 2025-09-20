package tui

import (
	"log"

	"github.com/qzeleza/termos"
)

// SelectMainMenu выбирает режим работы приложения
func (ac *AppConfig) SelectMainMenu() {

	// Устанавливаем язык для Термоса (TUI)
	termos.SetDefaultLanguage("en")
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Библиотека приложений для роутера").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Выводим информацию о системе
	ac.SysInfo(setupQueue)

	// Создаем список для выбора
	menuList := []string{
		"Приложения",
		"Настройки",
		"Выход",
	}

	// Создаем задачу для выбора пункта меню с запоминанием последней позиции
	menuTask := termos.NewSingleSelectTask("Выбор режима работы", menuList).WithDefaultItem(ac.LastMainMenuIndex)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		log.Fatal("Ошибка при выборе режима работы:", err)
	}
	// Сохраняем выбранный индекс и устанавливаем режим
	ac.LastMainMenuIndex = menuTask.GetSelectedIndex()
	ac.Mode = menuList[menuTask.GetSelectedIndex()]
}

//
