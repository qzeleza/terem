package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/qzeleza/terem/cmd/args"
	"github.com/qzeleza/terem/cmd/tui"
)

func main() {

	DEBUG := false
	VERSION := "1.0.0"
	APPNAME := "terem"
	LOGFILE := fmt.Sprintf("/tmp/%s.log", APPNAME)
	CONF := fmt.Sprintf("/opt/etc/%s/config.yaml", APPNAME)

	// 1. Инициализируем конфигурацию приложени
	ac, err := tui.NewSetup(APPNAME, VERSION, DEBUG)
	if err != nil {
		fmt.Printf("Ошибка создания конфигурации: %v\n", err)
		os.Exit(1)
	}

	// 2. Устанавливаем файлы логов и конфигурации
	ac.LogFile = LOGFILE
	ac.ConfFile = CONF

	// 3. Восстановление паники в случае ошибки
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

	// 4. Включаем обработчик SIGHUP
	defer ac.Log.Close() // Закрываем логгер при завершении программы
	select {} // приложение продолжает работать, ожидая SIGHUP
}

	// 4. Создаем корневой контекст для graceful shutdown
	var rootCancel context.CancelFunc
	ac.RootCtx, rootCancel = context.WithCancel(context.Background())
	defer rootCancel()

	// 5. Запускаем приложение c обработкой аргументов командной строки
	args.Execute(ac)
}
