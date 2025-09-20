package tui

import (
	"log"

	"github.com/qzeleza/terem/internal/i18n"
	"github.com/qzeleza/termos"
)

// otherList содержит список прочих приложений в фиксированном порядке
var otherList = []string{
	OtherOptionInfo,
	OtherOptionBack,
}

// otherKeys соответствующие ключи для otherList
var otherKeys = []string{
	"info",
	"quit",
}

// OtherCategoryLoop запускает цикл для выбора прочих приложений
func (ac *AppConfig) OtherCategoryLoop() {
	ac.ContextualLoop(func() bool {
		ac.SelectOtherCategory()

		// Проверяем контекст после выбора
		if ac.IsContextCancelled() {
			return false
		}

		switch ac.Category {
		case OtherOptionInfo:
			ac.SelectInfoApp()
			return true
		case OtherOptionBack:
			return false
		default:
			ac.Log.Warn(i18n.T("others.warn.invalid"))
			return false
		}
	}, i18n.T("loop.others"))
}

// SelectOtherCategory отображает меню для выбора прочих приложений
func (ac *AppConfig) SelectOtherCategory() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue(i18n.T("others.queue.title")).
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	labels := labelsFor(otherList)

	// Создаем задачу для выбора пункта меню с запоминанием последней позиции
	menuTask := termos.NewSingleSelectTask(i18n.T("others.task.title"), labels).WithDefaultItem(ac.LastOthersIndex)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		log.Fatal(i18n.T("others.error"), err)
	}

	// Сохраняем выбранный индекс и устанавливаем категорию
	selected := menuTask.GetSelectedIndex()
	ac.LastOthersIndex = selected
	ac.Category = otherList[selected]
}

// SelectInfoApp отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectInfoApp() {
	ac.Log.Info(i18n.T("others.log.info"))
}
