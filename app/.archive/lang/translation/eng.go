package translation

// eng - словарь для английского языка.
// Ключ - фраза на русском, значение - перевод на английский.
var Eng = map[string]string{
	// Системные сообщения
	"ПАНИКА": "PANIC",
	"Ошибка создания логгера": "Logger creation error",

	// Команды и справка
	"Терем - утилита для управления роутерами": "Terem - router management utility",
	"Терем - это утилита для упрощения работы с утилитами на роутерах\nс entware/openwrt. Поддерживает интерактивный режим и команды.\n\nИспользование:\n  terem           - запуск в интерактивном режиме\n  terem info      - информация о системе\n  terem [command] - выполнение конкретной команды": "Terem is a utility for simplifying work with utilities on routers\nwith entware/openwrt. Supports interactive mode and commands.\n\nUsage:\n  terem           - launch in interactive mode\n  terem info      - system information\n  terem [command] - execute specific command",
	"Для справки используйте: terem --help":                    "For help use: terem --help",
	"Для запуска интерактивного режима используйте: terem app": "To launch interactive mode use: terem app",

	// Команда info
	"Информация о системе": "System information",
	"Отображает информацию о системе в полном объеме:\nверсии программного обеспечения, аппаратные характеристики и т.д.": "Displays complete system information:\nsoftware versions, hardware specifications, etc.",
	"=== Информация о системе ===": "=== System Information ===",
	"Go версия:": "Go version:",
	"Архитектура: ARM (по умолчанию для роутеров)": "Architecture: ARM (default for routers)",

	// Команда app
	"Запуск основного приложения Терем": "Launch main Terem application",
	"Запускает основное приложение Терем для управления утилитами\nна роутерах с entware/openwrt. Если команда запущена без аргументов,\nбудет запущен интерактивный режим.": "Launches the main Terem application for managing utilities\non routers with entware/openwrt. If run without arguments,\ninteractive mode will be launched.",
	"Инициализация приложения...":          "Initializing application...",
	"Запуск основного приложения Терем...": "Starting main Terem application...",
	"TUI интерфейс будет реализован позже": "TUI interface will be implemented later",
	"Завершение работы приложения...":      "Shutting down application...",
	"Ошибка запуска приложения:":           "Application startup error:",
}
