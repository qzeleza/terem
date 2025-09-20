package args

import (
	"fmt"

	"github.com/qzeleza/terem/internal/i18n"
	"github.com/spf13/cobra"
)

// infoCmd команда для отображения информации о системе
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: i18n.T("cli.info.short"),
	Long:  i18n.T("cli.info.long"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(i18n.T("cli.info.header"))
		fmt.Println(i18n.T("cli.info.version"))
		fmt.Println(fmt.Sprintf(i18n.T("cli.info.go_version"), "1.25.0+"))
		fmt.Println(fmt.Sprintf(i18n.T("cli.info.arch"), i18n.T("cli.info.arch.value")))
	},
}

func localizeInfoCommand() {
	infoCmd.Short = i18n.T("cli.info.short")
	infoCmd.Long = i18n.T("cli.info.long")
}

func init() {
	localizeInfoCommand()
	// Добавляем команду info
	rootCmd.AddCommand(infoCmd)
}
