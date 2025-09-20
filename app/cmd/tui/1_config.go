package tui

import (
	"context"
	"fmt"
	"sync"

	"github.com/charmbracelet/lipgloss"
	conf "github.com/qzeleza/terem/internal/config"
	"github.com/qzeleza/terem/internal/i18n"
	log "github.com/qzeleza/terem/internal/zlog"
	"github.com/qzeleza/termos"
)

type SelectedApp struct {
	Name        string
	Description string
}

// Структура для хранения основных компонентов приложения
type AppConfig struct {
	AppName       string
	AppTitleColor lipgloss.TerminalColor
	AppTitle      string
	LogFile       string
	ConfFile      string
	Conf          conf.Config
	Log           log.Logger
	RootCtx       context.Context
	CancelFunc    context.CancelFunc // Функция для отмены контекста при shutdown
	Version       string
	Mode          string
	Category      string
	Debug         bool
	Language      string
	SelectedUtil  SelectedApp
	// Поля для кеширования системной информации
	cachedSysInfo *SysInfoResult
	sysInfoOnce   sync.Once
	// Поля для запоминания последних выбранных позиций в меню
	LastMainMenuIndex int // Главное меню (Приложения/Настройки/Выход)
	LastCategoryIndex int // Меню категорий (Безопасность/Сетевые/Прочие/Назад)
	LastSecurityIndex int // Меню безопасности
	LastNetworkIndex  int // Меню сетевых утилит
	LastOthersIndex   int // Меню прочих утилит
	LastSettingsIndex int // Меню настроек
}

func NewSetup(language string, appName string, version string, debug bool, logFile string, confFile string) (*AppConfig, error) {
	// Загружаем конфигурацию
	confData, resolvedPath, err := conf.Load(confFile)
	if err != nil {
		return nil, err
	}

	// Определяем путь до файла логов
	if confData.LogFile != "" {
		logFile = confData.LogFile
	}

	// Устанавливаем язык
	if err := i18n.SetLanguage(language); err != nil {
		fmt.Printf(i18n.T("language.warn.unsupported")+"\n", language)
		language = "ru"
		_ = i18n.SetLanguage(language)
	}

	if err := i18n.Error(); err != nil {
		return nil, err
	}

	confData.SetLanguage(language)

	ac := &AppConfig{
		AppName:       appName,
		AppTitleColor: termos.GreenBright,
		AppTitle:      i18n.T("app.title"),
		LogFile:       logFile,
		ConfFile:      resolvedPath,
		Conf:          *confData,
		Version:       version,
		Debug:         debug,
		Language:      i18n.Language(),
		SelectedUtil: SelectedApp{
			Name:        "",
			Description: "",
		},
	}

	// Инициализируем логгер
	if err := ac.SetupLogger(); err != nil {
		return nil, err
	}

	return ac, nil
}
func (ac *AppConfig) SetupLogger() error {
	logger := log.New(ac.LogFile)

	if ac.Conf.DebugMode || ac.Debug {
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}
	ac.Log = *logger

	return nil
}

// IsContextCancelled проверяет, был ли отменен контекст (например, по Ctrl+C)
func (ac *AppConfig) IsContextCancelled() bool {
	if ac.RootCtx == nil {
		return false
	}
	select {
	case <-ac.RootCtx.Done():
		return true
	default:
		return false
	}
}

// GracefulShutdown выполняет graceful shutdown приложения
func (ac *AppConfig) GracefulShutdown() {
	ac.Log.Info(i18n.T("shutdown.log.start"))
	if ac.CancelFunc != nil {
		ac.CancelFunc()
	}
}

// ContextualLoop - обертка для циклов с автоматической проверкой контекста
// loopBody должна возвращать true для продолжения цикла, false для выхода
func (ac *AppConfig) ContextualLoop(loopBody func() bool, loopName string) {
	for {
		// Проверяем контекст перед каждой итерацией
		if ac.IsContextCancelled() {
			ac.Log.Info(fmt.Sprintf(i18n.T("loop.signal"), loopName))
			return
		}

		// Выполняем тело цикла
		if !loopBody() {
			return // выход из цикла по решению loopBody
		}
	}
}

// GetSysInfo возвращает кешированную системную информацию
// Использует sync.Once для гарантии однократного выполнения
func (ac *AppConfig) GetSysInfo() *SysInfoResult {
	ac.sysInfoOnce.Do(func() {
		ac.Log.Debug(i18n.T("sysinfo.log.first"))
		ac.cachedSysInfo = &SysInfoResult{}
		ac.getSysInfo(ac.cachedSysInfo)
		ac.Log.Debug(i18n.T("sysinfo.log.cache"))
	})
	return ac.cachedSysInfo
}
