package terem

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	// tms "github.com/qzeleza/termos"
	conf "github.com/qzeleza/terem/internal/config"
	log "github.com/qzeleza/zlogger"
)

// Имя основного приложения
const AppName = "terem"

// Структура для хранения основных компонентов приложения
type appConfig struct {
	logFile  string
	sockFile string
	conf     conf.Config
	logger   log.Logger
	rootCtx  context.Context
}

func main() {

	// 1. Восстановление паники
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ПАНИКА: %v\n", r)
			// Выводим стек вызовов
			debug.PrintStack()
		}
	}()

	// 3. Задаем параметры конфигураци
	appConf := appConfig{
		logFile:  fmt.Sprintf("/tmp/%s.log", AppName),
		sockFile: fmt.Sprintf("/tmp/%s.sock", AppName),
		conf:     *conf.New(),
	}

	// 2. Создаем корневой контекст для graceful shutdown
	var rootCancel context.CancelFunc
	appConf.rootCtx, rootCancel = context.WithCancel(context.Background())
	defer rootCancel()

	// 4. Инициализируем логгер приложения
	err := initLogger(&appConf)
	if err != nil {
		fmt.Printf("Ошибка создания логгера: %v\n", err)
		os.Exit(1)
	}
	defer appConf.logger.Close()

}
