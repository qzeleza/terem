package tui

import (
	"log"

	"github.com/qzeleza/termos"
)

// otherList содержит список прочих приложений в фиксированном порядке
var otherList = []string{
	"Информация о системе",
	"Назад",
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
		case otherList[0]: // Информация о системе
			ac.SelectInfoApp()
			// После выполнения действия показываем меню снова
			return true
		case otherList[1]: // Назад
			return false
		default:
			ac.Log.Warn("Неверный выбор категории")
			return false
		}
	}, "цикла прочих приложений")
}

// SelectOtherCategory отображает меню для выбора прочих приложений
func (ac *AppConfig) SelectOtherCategory() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор прочих приложений").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем задачу для выбора пункта меню с запоминанием последней позиции
	menuTask := termos.NewSingleSelectTask("Выбор приложения", otherList).WithDefaultItem(ac.LastOthersIndex)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		log.Fatal("Ошибка при выборе приложения:", err)
	}

	// Сохраняем выбранный индекс и устанавливаем категорию
	ac.LastOthersIndex = menuTask.GetSelectedIndex()
	ac.Category = otherList[menuTask.GetSelectedIndex()]
}

// SelectInfoApp отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectInfoApp() {
	ac.Log.Info("Выбрано приложение для информации о системе")
}
