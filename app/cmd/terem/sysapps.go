package terem

import (
	"github.com/qzeleza/termos"
)

var fsAppList = []string{
	"Резервное копирование конфигурации",
	"Выход",
}

func (ac *AppConfig) SelectFileSystemApp() {
	for {
		ac.SelectFileSystemAppFromList()
		switch ac.Mode {
		case fsAppList[1]: // Выход
			return
		}
	}
}

// SelectFileSystemApp отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectFileSystemAppFromList() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор утилиты для работы с файловой системой").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем задачу для выбора пункта меню
	menuTask := termos.NewSingleSelectTask("Выбор утилиты", fsAppList)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		ac.Log.Fatal("Ошибка при выборе утилиты:", err)
	}

	ac.Mode = fsAppList[menuTask.GetSelectedIndex()]
}
