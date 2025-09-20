package tui

import (
	"github.com/qzeleza/terem/internal/i18n"
	"github.com/qzeleza/termos"
)

var categoryList = []string{
	CategorySecurity,
	CategoryNetwork,
	CategoryOther,
	CategoryBack,
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
		case CategorySecurity:
			ac.SecurityCategoryLoop()
		case CategoryNetwork:
			ac.NetworkCategoryLoop()
		case CategoryOther:
			ac.OtherCategoryLoop()
		case CategoryBack:
			return false
		default:
			ac.Log.Warn(i18n.T("category.warn.invalid"))
			return false
		}

		// После возврата из подменю показываем меню категорий снова
		return true
	}, i18n.T("loop.category"))
}

// SelectCategory отображает меню для выбора категории приложений
func (ac *AppConfig) SelectCategoryFromList() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue(i18n.T("category.queue.title")).
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	labels := labelsFor(categoryList)

	// Создаем задачу для выбора пункта меню с запоминанием последней позиции
	menuTask := termos.NewSingleSelectTask(i18n.T("category.task.title"), labels).WithDefaultItem(ac.LastCategoryIndex)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		ac.Log.Fatal(i18n.T("category.error"), err)
	}

	// Сохраняем выбранный индекс и устанавливаем категорию
	selected := menuTask.GetSelectedIndex()
	ac.LastCategoryIndex = selected
	ac.Category = categoryList[selected]
}
