package tui

import (
	"github.com/qzeleza/terem/internal/i18n"
	"github.com/qzeleza/termos"
)

// Список меню для выбора категорий безопасности в фиксированном порядке
var securityAppList = []string{
	SecurityOptionParental,
	SecurityOptionAntiscan,
	SecurityOptionBackup,
	SecurityOptionBack,
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
		case SecurityOptionParental:
			ac.SelectParentalControl()
			return true
		case SecurityOptionAntiscan:
			ac.SelectAntiscan()
			return true
		case SecurityOptionBackup:
			ac.SelectBackup()
			return true
		case SecurityOptionBack:
			return false
		default:
			ac.Log.Warn(i18n.T("security.warn.invalid"))
			return false
		}
	}, i18n.T("loop.security"))
}

// SelectSecurityApp отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectSecurityApp() {
	// Создаем основную очередь для выбора приложения
	setupQueue := termos.NewQueue(i18n.T("security.queue.title")).
		WithAppName(ac.AppTitle).
		WithSummary(false).
		WithTitleColor(ac.AppTitleColor, true).
		WithClearScreen(true)

	labels := labelsFor(securityAppList)

	// Создаем задачу для выбора пункта меню с запоминанием последней позиции
	menuTask := termos.NewSingleSelectTask(i18n.T("security.task.title"), labels).WithDefaultItem(ac.LastSecurityIndex)
	setupQueue.AddTasks(menuTask)

	// Запускаем выбор режима
	if err := setupQueue.Run(); err != nil {
		ac.Log.Fatal(i18n.T("security.error"), err)
	}

	// Если пользователь отменил выбор, возвращаемся в предыдущее меню
	if menuTask.HasError() {
		ac.Mode = SecurityOptionBack
		return
	}

	// Сохраняем выбранный индекс и устанавливаем режим
	selected := menuTask.GetSelectedIndex()
	ac.LastSecurityIndex = selected
	ac.Mode = securityAppList[selected]
}

// SelectParentalControl отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectParentalControl() {
	ac.Log.Info(i18n.T("security.log.parental"))
}

// SelectAntiscan отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectAntiscan() {
	ac.Log.Info(i18n.T("security.log.antiscan"))
}

// SelectAntiscan отображает меню для выбора утилит для работы с файловой системой
func (ac *AppConfig) SelectBackup() {
	ac.Log.Info(i18n.T("security.log.backup"))
}
