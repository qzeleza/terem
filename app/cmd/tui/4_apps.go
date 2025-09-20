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
	ac.ContextualLoop(func() bool {
		ac.SelectCategoryFromList()

		// Проверяем контекст после выбора
		if ac.IsContextCancelled() {
			return false
		}

		switch ac.Category {
		case categoryList[0]: // Безопасность роутера
			ac.SecurityCategoryLoop()
		case categoryList[1]: // Сетевые утилиты
			ac.NetworkCategoryLoop()
		case categoryList[2]: // Прочие утилиты
			ac.OtherCategoryLoop()
		case categoryList[3]: // Выход
			return false
		default:
			ac.Log.Warn("Неверный выбор категории")
			return false
		}

		// После возврата из подменю показываем меню категорий снова
		return true
	}, "выбора категории приложений")
}

// SelectCategory отображает меню для выбора категории приложений
func (ac *AppConfig) SelectCategoryFromList() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор категории приложений").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем задачу для выбора пункта меню с запоминанием последней позиции
	menuTask := termos.NewSingleSelectTask("Выбор категории", categoryList).WithDefaultItem(ac.LastCategoryIndex)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		ac.Log.Fatal("Ошибка при выборе категории:", err)
	}

	// Сохраняем выбранный индекс и устанавливаем категорию
	ac.LastCategoryIndex = menuTask.GetSelectedIndex()
	ac.Category = categoryList[menuTask.GetSelectedIndex()]
}
