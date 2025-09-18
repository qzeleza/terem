package terem

import (
	"log"

	"github.com/qzeleza/termos"
)

var networkList = []map[string]string{
	{"nfs": "Сервер NFS"},
	{"exit": "Выход"},
}

func (ac *AppConfig) SelectNetworkApp() {
	for {
		ac.SelectNetworkAppFromList()
		switch ac.Category {
		case networkList[0]["nfs"]: // Сервер NFS
			ac.SelectFileSystemApp()
		case networkList[1]["exit"]: // Выход
			return
		default:
			ac.Log.Warn("Неверный выбор категории")
			return
		}
	}
}

func (ac *AppConfig) SelectNetworkAppFromList() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор сетевых приложений").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем список для выбора
	list := []string{}
	for _, v := range networkList {
		for _, v := range v {
			list = append(list, v)
		}
	}

	// Создаем задачу для выбора пункта меню
	menuTask := termos.NewSingleSelectTask("Выберите сетевое приложение", list)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		log.Fatal("Ошибка при выборе сетевого приложения:", err)
	}

	ac.SelectedUtil.Description = list[menuTask.GetSelectedIndex()]
}
