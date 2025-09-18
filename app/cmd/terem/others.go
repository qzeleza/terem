package terem

import (
	"log"

	"github.com/qzeleza/termos"
)

var otherList = []map[string]string{
	{"dns": "DNS-сервер"},
	{"exit": "Выход"},
}

func (ac *AppConfig) SelectOtherApp() {
	for {
		ac.SelectOtherAppFromList()
		switch ac.Category {
		case otherList[0]["dns"]: // DNS-сервер
			ac.SelectFileSystemApp()
		case otherList[1]["exit"]: // Выход
			return
		default:
			ac.Log.Warn("Неверный выбор категории")
			return
		}
	}
}

func (ac *AppConfig) SelectOtherAppFromList() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор прочих приложений").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем список для выбора
	list := []string{}
	for _, v := range otherList {
		for _, v := range v {
			list = append(list, v)
		}
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
