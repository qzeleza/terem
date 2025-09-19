package tui

import (
	"log"

	"github.com/qzeleza/termos"
)

// otherList содержит список прочих приложений
var otherList = map[string]string{
	"info": "Информация о системе",
	"quit": "Назад",
}

// OtherCategoryLoop запускает цикл для выбора прочих приложений
func (ac *AppConfig) OtherCategoryLoop() {
	for {
		ac.SelectOtherCategory()
		switch ac.Category {
		case otherList["info"]: // Информация о системе
			ac.SelectInfoApp()
		case otherList["exit"]: // Выход
			return
		default:
			ac.Log.Warn("Неверный выбор категории")
			return
		}
	}
}

// SelectOtherCategory отображает меню для выбора прочих приложений
func (ac *AppConfig) SelectOtherCategory() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор прочих приложений").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем список для выбора
	list := []string{}
	for _, v := range otherList {
		list = append(list, v)
	}

	// Создаем задачу для выбора пункта меню
	menuTask := termos.NewSingleSelectTask("Выбор приложения", list)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		log.Fatal("Ошибка при выборе приложения:", err)
	}

	ac.SelectedUtil.Description = list[menuTask.GetSelectedIndex()]
}

// SelectInfoApp отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectInfoApp() {
	ac.Log.Info("Выбрано приложение для информации о системе")
}
