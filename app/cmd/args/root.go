package args

import (
	"fmt"
	"os"

	"github.com/qzeleza/terem/cmd/tui"
	"github.com/spf13/cobra"
)

// // Интерфейс для конфигурации приложения
// type AppConfigInterface interface {
// 	SelectMode() // Выбор режима работы
// }

// // Глобальные переменные для флагов
// // Будет передан из main
var AppConfig *tui.AppConfig

var rootCmd = &cobra.Command{
	Use:   "terem",
	Short: "Терем - утилита для управления роутерами",
	Long: `Терем - это утилита для упрощения работы с утилитами на роутерах
				с entware/openwrt. Поддерживает интерактивный режим и команды.

				Использование:
				terem           - запуск в интерактивном режиме
				terem info      - информация о системе
				terem [command] - выполнение конкретной команды`,
	Run: func(cmd *cobra.Command, args []string) {
		// Запускаем интерактивный режим, в случае если запущена без аргументов
		if len(args) == 0 {

			// Запускаем главный цикл с автоматической проверкой контекста
			AppConfig.ContextualLoop(func() bool {
				AppConfig.SelectMainMenu()

				// Проверяем контекст после выбора меню
				if AppConfig.IsContextCancelled() {
					return false
				}

				switch AppConfig.Mode {
				case "Приложения":
					AppConfig.SelectCategoryLoop()
				case "Настройки":
					AppConfig.SelectSettingsLoop()
				case "Выход":
					AppConfig.Log.Info("Пользователь выбрал выход")
					return false
				}

				return true // продолжить главный цикл
			}, "главного меню")
		}
	},
}

// Execute запускает командную строку
func Execute(ac *tui.AppConfig) {
	ac.Log.Info("Запуск командной строки")
	AppConfig = ac
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
