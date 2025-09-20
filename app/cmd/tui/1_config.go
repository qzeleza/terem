package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	conf "github.com/qzeleza/terem/internal/config"
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
	Version       string
	Mode          string
	Category      string
	Debug         bool
	SelectedUtil  SelectedApp
}

func NewSetup(appName string, version string, debug bool) (*AppConfig, error) {
	// Загружаем конфигурацию
	defaultConfPath := fmt.Sprintf("/opt/etc/%s/config.yaml", appName)
	confData, resolvedPath, err := conf.Load(defaultConfPath)
	if err != nil {
		return nil, err
	}

	// Определяем путь до файла логов
	logFile := fmt.Sprintf("/tmp/%s.log", appName)
	if confData.LogFile != "" {
		logFile = confData.LogFile
	}

	ac := &AppConfig{
		AppName:       appName,
		AppTitleColor: termos.GreenBright,
		AppTitle:      "Терем™",
		LogFile:       logFile,
		ConfFile:      resolvedPath,
		Conf:          *confData,
		Version:       version,
		Debug:         debug,
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
