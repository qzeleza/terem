package args

import (
	"github.com/spf13/cobra"
)

// debugCmd команда для отображения информации о системе
var debugCmd = &cobra.Command{
	Use:     "--debug",
	Aliases: []string{"d", "debug"},
	Short:   "Отображение отладочной информации",
	Long: `Отображает отладочную информацию пакета в файле логов:
			версии программного обеспечения, аппаратные характеристики и т.д.`,
	Run: func(cmd *cobra.Command, args []string) {
		AppConfig.Debug = true
		err := AppConfig.SetupLogger()
		if err != nil {
			err := AppConfig.Log.Fatal("Ошибка при настройке логгера:", err)
			if err != nil {
				AppConfig.Log.Error("Ошибка при выводе сообщения:", err)
			}
		}
	},
}

func init() {
	// Добавляем команду debug
	rootCmd.AddCommand(debugCmd)
}
