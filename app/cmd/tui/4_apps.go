package tui

import (
	"github.com/qzeleza/termos"
)

var categoryList = []string{
	"Безопасность роутера",
	"Сетевые утилиты",
	"Прочие утилиты",
	"Назад",
}

// SelectCategory отображает меню для выбора категории приложений
func (ac *AppConfig) SelectCategoryLoop() {
	for {
		ac.SelectCategoryFromList()
		switch ac.Category {
		case categoryList[0]: // Безопасность роутера
			ac.SecurityCategoryLoop()
		case categoryList[1]: // Сетевые утилиты
			ac.NetworkCategoryLoop()
		case categoryList[2]: // Прочие утилиты
			ac.OtherCategoryLoop()
		case categoryList[3]: // Выход
			return
		default:
			ac.Log.Warn("Неверный выбор категории")
			return
		}
	}
}

// SelectCategory отображает меню для выбора категории приложений
func (ac *AppConfig) SelectCategoryFromList() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор категории приложений").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем задачу для выбора пункта меню
	menuTask := termos.NewSingleSelectTask("Выбор категории", categoryList)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		ac.Log.Fatal("Ошибка при выборе категории:", err)
	}

	ac.Category = categoryList[menuTask.GetSelectedIndex()]
}
