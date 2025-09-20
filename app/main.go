package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/qzeleza/terem/cmd/args"
	"github.com/qzeleza/terem/cmd/tui"
)

func main() {

	DEBUG := false
	VERSION := "1.0.0"
	APPNAME := "terem"
	LOGFILE := fmt.Sprintf("/tmp/%s.log", APPNAME)
	CONF := fmt.Sprintf("/opt/etc/%s/config.yaml", APPNAME)

	// 1. Инициализируем конфигурацию приложения
	ac, err := tui.NewSetup(APPNAME, VERSION, DEBUG, LOGFILE, CONF)
	if err != nil {
		fmt.Printf("Ошибка создания конфигурации: %v\n", err)
		os.Exit(1)
	}

	// 2. Устанавливаем файлы логов и конфигурации
	ac.LogFile = LOGFILE
	ac.ConfFile = CONF

	// 3. Создаем корневой контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	ac.RootCtx = ctx
	ac.CancelFunc = cancel
	defer cancel()

	// 4. Настраиваем обработку сигналов для graceful shutdown
	setupSignalHandler(cancel, ac)

	// 5. Восстановление паники в случае ошибки
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("ПАНИКА: %v\n%s", r, debug.Stack())
			fmt.Fprintf(os.Stderr, "%s\n", msg)

			// Если логгер уже инициализирован, логируем ошибку
			if ac != nil {
				ac.Log.Error("Критическая ошибка: ", msg)
			}

			os.Exit(1)
		}
	}()

	// 6. Закрываем логгер при завершении программы
	defer ac.Log.Close()

	// 7. Запускаем приложение c обработкой аргументов командной строки
	args.Execute(ac)
}

// setupSignalHandler настраивает обработку сигналов для graceful shutdown
func setupSignalHandler(cancel context.CancelFunc, ac *tui.AppConfig) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		ac.Log.Info(fmt.Sprintf("Получен сигнал %v, начинаем graceful shutdown", sig))
		cancel()
	}()
}
