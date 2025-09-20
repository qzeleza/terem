package tui

import "github.com/qzeleza/terem/internal/i18n"

const (
	ModeApps     = "menu.main.option.apps"
	ModeSettings = "menu.main.option.settings"
	ModeExit     = "menu.main.option.exit"

	CategorySecurity = "category.security"
	CategoryNetwork  = "category.network"
	CategoryOther    = "category.other"
	CategoryBack     = "category.back"

	NetworkOptionOpenSSH = "network.option.openssh"
	NetworkOptionProxy   = "network.option.proxy"
	NetworkOptionDNS     = "network.option.dns"
	NetworkOptionAdGuard = "network.option.adguard"
	NetworkOptionBack    = "network.option.back"

	OtherOptionInfo = "others.option.info"
	OtherOptionBack = "others.option.back"

	SecurityOptionParental = "security.option.parental"
	SecurityOptionAntiscan = "security.option.antiscan"
	SecurityOptionBackup   = "security.option.backup"
	SecurityOptionBack     = "security.option.back"

	SettingsOptionLogging = "settings.option.logging"
	SettingsOptionBack    = "settings.option.back"
)

func labelsFor(keys []string) []string {
	items := make([]string, len(keys))
	for i, key := range keys {
		items[i] = i18n.T(key)
	}
	return items
}
