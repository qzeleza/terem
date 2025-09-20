package tui

import (
	"github.com/qzeleza/termos"
)

// Список меню для выбора категорий безопасности в фиксированном порядке
var securityAppList = []string{
	"Родительский контроль",
	"Защита роутера от атак Antiscan",
	"Резервное копирование конфигурации",
	"Назад",
}

// securityKeys соответствующие ключи для securityAppList
var securityKeys = []string{
	"parental",
	"antiscan",
	"backup",
	"quit",
}

// SecurityCategoryLoop запускает цикл для выбора категорий безопасности
func (ac *AppConfig) SecurityCategoryLoop() {
	ac.ContextualLoop(func() bool {
		ac.SelectSecurityApp()

		// Проверяем контекст после выбора
		if ac.IsContextCancelled() {
			return false
		}

		switch ac.Mode {
		case securityAppList[0]: // Родительский контроль
			ac.SelectParentalControl()
			// После выполнения действия показываем меню снова
			return true
		case securityAppList[1]: // Защита роутера от атак Antiscan
			ac.SelectAntiscan()
			// После выполнения действия показываем меню снова
			return true
		case securityAppList[2]: // Резервное копирование конфигурации
			ac.SelectBackup()
			// После выполнения действия показываем меню снова
			return true
		case securityAppList[3]: // Назад
			return false
		default:
			ac.Log.Warn("Неверный выбор приложения безопасности")
			return false
		}
	}, "цикла безопасности")
}

// SelectSecurityApp отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectSecurityApp() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue("Выбор программ для безопасности роутера").
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	// Создаем задачу для выбора пункта меню с запоминанием последней позиции
	menuTask := termos.NewSingleSelectTask("Выбор утилиты", securityAppList).WithDefaultItem(ac.LastSecurityIndex)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		ac.Log.Fatal("Ошибка при выборе утилиты:", err)
	}

	// Сохраняем выбранный индекс и устанавливаем режим
	ac.LastSecurityIndex = menuTask.GetSelectedIndex()
	ac.Mode = securityAppList[menuTask.GetSelectedIndex()]
}

// SelectParentalControl отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectParentalControl() {
	ac.Log.Info("Выбран родительский контроль")
}

// SelectAntiscan отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectAntiscan() {
	ac.Log.Info("Выбрана защита роутера от атак Antiscan")
}

// SelectAntiscan отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectBackup() {
	ac.Log.Info("Выбрано резервное копирование конфигурации")
}
