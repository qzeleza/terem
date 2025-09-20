package args

import (
	"github.com/qzeleza/terem/internal/i18n"
	"github.com/spf13/cobra"
)

// debugCmd команда для отображения информации о системе
var debugCmd = &cobra.Command{
	Use:     "--debug",
	Aliases: []string{"d", "debug"},
	Short:   i18n.T("cli.debug.short"),
	Long:    i18n.T("cli.debug.long"),
	Run: func(cmd *cobra.Command, args []string) {
		AppConfig.Debug = true
		err := AppConfig.SetupLogger()
		if err != nil {
			err := AppConfig.Log.Fatal(i18n.T("cli.debug.error"), err)
			if err != nil {
				AppConfig.Log.Error(i18n.T("cli.debug.write_error"), err)
			}
		}
	},
}

func localizeDebugCommand() {
	debugCmd.Short = i18n.T("cli.debug.short")
	debugCmd.Long = i18n.T("cli.debug.long")
}

func init() {
	localizeDebugCommand()
	// Добавляем команду debug
	rootCmd.AddCommand(debugCmd)
}
