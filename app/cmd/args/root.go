package args

import (
	"fmt"
	"os"

	"github.com/qzeleza/terem/cmd/tui"
	"github.com/qzeleza/terem/internal/i18n"
	"github.com/spf13/cobra"
)

// // Интерфейс для конфигурации приложения
// type AppConfigInterface interface {
// 	SelectMode() // Выбор режима работы
// }

// // Глобальные переменные для флагов
// // Будет передан из main
var AppConfig *tui.AppConfig
var languageFlag string

// rootCmd - основная команда
var rootCmd = &cobra.Command{
	Use:   "terem",
	Short: i18n.T("cli.root.short"),
	Long:  i18n.T("cli.root.long"),
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
				case tui.ModeApps:
					AppConfig.SelectCategoryLoop()
				case tui.ModeSettings:
					AppConfig.SelectSettingsLoop()
				case tui.ModeExit:
					AppConfig.Log.Info(i18n.T("menu.main.log.exit"))
					return false
				}

				return true // продолжить главный цикл
			}, i18n.T("loop.main"))
		}
	},
}

// локализация команды root
func localizeRoot() {
	rootCmd.Short = i18n.T("cli.root.short")
	rootCmd.Long = i18n.T("cli.root.long")
	localizeNetworkCommand()
	localizeDebugCommand()
	localizeInfoCommand()
}

func applyLanguageOverride() {
	if languageFlag == "" {
		return
	}

	if err := i18n.SetLanguage(languageFlag); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
		return
	}

	if AppConfig != nil {
		AppConfig.Language = i18n.Language()
		AppConfig.Conf.SetLanguage(AppConfig.Language)
		AppConfig.AppTitle = i18n.T("app.title")
	}

	localizeRoot()
}

// Execute запускает командную строку
func Execute(ac *tui.AppConfig) {
	ac.Log.Info(i18n.T("cli.root.log.start"))
	AppConfig = ac // передаем конфигурацию
	localizeRoot() // локализация команды root
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(applyLanguageOverride)
	rootCmd.PersistentFlags().StringVarP(&languageFlag, "lang", "l", "", "interface language (ru, en, tt)")
}
