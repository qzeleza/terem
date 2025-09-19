package tui

import (
	"github.com/qzeleza/termos"
)

// Список меню для выбора категорий безопасности
var securityAppList = map[string]string{
	"parental": "Родительский контроль",
	"antiscan": "Защита роутера от атак Antiscan",
	"backup":   "Резервное копирование конфигурации",
	"quit":     "Назад",
}

// SecurityCategoryLoop запускает цикл для выбора категорий безопасности
func (ac *AppConfig) SecurityCategoryLoop() {
	for {
		ac.SelectSecurityApp()
		switch ac.Mode {
		case securityAppList["parental"]: // Родительский контроль
			ac.SelectParentalControl()
		case securityAppList["antiscan"]: // Защита роутера от атак Antiscan
			ac.SelectAntiscan()
		case securityAppList["backup"]: // Резервное копирование конфигурации
			ac.SelectBackup()
		case securityAppList["exit"]: // Назад
			return
		}
	}
}

// SelectSecurityApp отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectSecurityApp() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор программ для безопасности роутера").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем список для выбора
	list := []string{}
	for _, v := range securityAppList {
		list = append(list, v)
	}

	// Создаем задачу для выбора пункта меню
	menuTask := termos.NewSingleSelectTask("Выбор утилиты", list)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		err := ac.Log.Fatal("Ошибка при выборе утилиты:", err)
		if err != nil {
			ac.Log.Error("Ошибка при выводе сообщения:", err)
		}
	}

	ac.Mode = list[menuTask.GetSelectedIndex()]
}

// SelectParentalControl отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectParentalControl() {
	err := ac.Log.Info("Выбран родительский контроль")
	if err != nil {
		ac.Log.Error("Ошибка при выводе сообщения:", err)
	}
}

// SelectAntiscan отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectAntiscan() {
	err := ac.Log.Info("Выбрана защита роутера от атак Antiscan")
	if err != nil {
		ac.Log.Error("Ошибка при выводе сообщения:", err)
	}
}

// SelectAntiscan отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectBackup() {
	err := ac.Log.Info("Выбрано резервное копирование конфигурации")
	if err != nil {
		ac.Log.Error("Ошибка при выводе сообщения:", err)
	}
}
