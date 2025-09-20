package args

import (
	"github.com/spf13/cobra"
)

// networkCmd команда для отображения сетевых приложений
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Отображает категории сетевых приложений",
	Long:  `Отображает все приложения из категории сетевых приложений для работы в интерактивном режиме`,
	Run: func(cmd *cobra.Command, args []string) {
		AppConfig.NetworkCategoryLoop()
	},
}

func init() {
	// Добавляем команду network
	rootCmd.AddCommand(networkCmd)
}
