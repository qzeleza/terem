package terem

import (
	"log"

	"github.com/qzeleza/termos"
)

// SelectMode выбирает режим работы приложения
func (ac *AppConfig) SelectMode() {

	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Библиотека приложений для роутера 'Терем'").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	menuList := []string{
		"Приложения",
		"Настройки",
		"Выход",
	}

	// Создаем задачу для выбора пункта меню
	menuTask := termos.NewSingleSelectTask("Выбор режима работы", menuList)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		log.Fatal("Ошибка при выборе режима работы:", err)
	}

	ac.Mode = menuList[menuTask.GetSelectedIndex()]
}
