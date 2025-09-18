package args

import (
	"fmt"

	"github.com/spf13/cobra"
)

// infoCmd команда для отображения информации о системе
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Информация о системе",
	Long: `Отображает информацию о системе в полном объеме:
версии программного обеспечения, аппаратные характеристики и т.д.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=== Информация о системе ===")
		fmt.Println("terem v1.0.0")
		fmt.Println("Go версия:", "1.25.0+")
		fmt.Println("Архитектура: ARM (по умолчанию для роутеров)")
	},
}

func init() {
	// Добавляем команду info
	rootCmd.AddCommand(infoCmd)
}
