package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	conf "github.com/qzeleza/terem/internal/config"
	"github.com/qzeleza/termos"
	log "github.com/qzeleza/zlogger"
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
	SockFile      string
	Conf          conf.Config
	Log           log.Logger
	RootCtx       context.Context
	Version       string
	Mode          string
	Category      string
	Debug         bool
	SelectedUtil  SelectedApp
}

func NewSetup(appName string, version string) (*AppConfig, error) {
	ac := &AppConfig{
		AppName:       appName,
		AppTitleColor: termos.GreenBright,
		AppTitle:      "Терем™",
		LogFile:       fmt.Sprintf("/tmp/%s.log", appName),
		SockFile:      fmt.Sprintf("/tmp/%s.sock", appName),
		Conf:          *conf.New(),
		Version:       version,
		Debug:         false,
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

	// Создаем конфигурацию с настройками по умолчанию
	logConfig := log.NewConfig(ac.LogFile, ac.SockFile)

	// Настраиваем дополнительные параметры
	if ac.Conf.DebugMode == "true" || ac.Debug {
		logConfig.Level = log.DEBUG.String()
	} else {
		logConfig.Level = log.INFO.String()
	}
	logConfig.MaxFileSize = 50                // максимальный размер файла логов
	logConfig.BufferSize = 500                // размер буфера
	logConfig.FlushInterval = 2 * time.Second // интервал обновления буфера

	// Создаем логгер для основного приложения
	logger, err := log.New(logConfig, "MAIN")
	if err != nil {
		return err
	}
	defer func() {
		if err := logger.Close(); err != nil {
			ac.Log.Error("Ошибка при закрытии логгера:", err)
		}
	}()

	ac.Log = *logger

	return nil
}
