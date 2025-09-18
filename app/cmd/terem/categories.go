package terem

import (
	"log"

	"github.com/qzeleza/termos"
)

var categoryList = []string{
	"Утилиты для работы с файловой системой",
	"Сетевые утилиты",
	"Прочие утилиты",
	"Выход",
}

func (ac *AppConfig) SelectCategory() {
	for {
		ac.SelectCategoryFromList()
		switch ac.Category {
		case categoryList[0]: // Утилиты для работы с файловой системой
			ac.SelectFileSystemApp()
		case categoryList[1]: // Сетевые утилиты
			ac.SelectNetworkApp()
		case categoryList[2]: // Прочие утилиты
			ac.SelectOtherApp()
		case categoryList[3]: // Выход
			return
		default:
			ac.Log.Warn("Неверный выбор категории")
			return
		}
	}
}

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
		log.Fatal("Ошибка при выборе категории:", err)
	}

	ac.Category = categoryList[menuTask.GetSelectedIndex()]
}
