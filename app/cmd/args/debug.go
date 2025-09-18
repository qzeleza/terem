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
		AppConfig.SetupLogger()
	},
}

func init() {
	// Добавляем команду debug
	rootCmd.AddCommand(debugCmd)
}
