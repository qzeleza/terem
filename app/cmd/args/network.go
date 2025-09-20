package args

import (
	"github.com/qzeleza/terem/internal/i18n"
	"github.com/spf13/cobra"
)

// networkCmd команда для отображения сетевых приложений
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: i18n.T("cli.network.short"),
	Long:  i18n.T("cli.network.long"),
	Run: func(cmd *cobra.Command, args []string) {
		AppConfig.NetworkCategoryLoop()
	},
}

func localizeNetworkCommand() {
	networkCmd.Short = i18n.T("cli.network.short")
	networkCmd.Long = i18n.T("cli.network.long")
}

func init() {
	localizeNetworkCommand()
	// Добавляем команду network
	rootCmd.AddCommand(networkCmd)
}
