package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/qzeleza/terem/cmd/args"
	"github.com/qzeleza/terem/cmd/terem"
)

func main() {

	// 1. Восстановление паники в случае ошибки
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ПАНИКА: %v\n", r)
			// Выводим стек вызовов
			debug.PrintStack()
		}
	}()

	// 2. Инициализируем конфигурацию приложени
	ac, err := terem.NewSetup("terem", "1.0.0")
	if err != nil {
		fmt.Printf("Ошибка создания конфигурации: %v\n", err)
		os.Exit(1)
	}

	// 3. Передаем конфигурацию в args для обработки флагов
	// args.SetAppConfig(ac)

	// 4. Создаем корневой контекст для graceful shutdown
	var rootCancel context.CancelFunc
	ac.RootCtx, rootCancel = context.WithCancel(context.Background())
	defer rootCancel()

	// 5. Запускаем приложение c обработкой аргументов командной строки
	args.Execute(ac)
}
