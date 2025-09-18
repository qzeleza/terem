package terem

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	// tms "github.com/qzeleza/termos"
	conf "github.com/qzeleza/terem/internal/config"
	"github.com/qzeleza/terem/internal/lang"
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
	lang     *lang.Translator
	version  string
}

// T является методом-помощником для appConfig для выполнения перевода.
func (ac *appConfig) T(rusString string) string {
	return ac.lang.T(rusString)
}

func main() {

	// 1. Задаем параметры конфигурации по умолчанию
	ac := appConfig{
		logFile:  fmt.Sprintf("/tmp/%s.log", AppName),
		sockFile: fmt.Sprintf("/tmp/%s.sock", AppName),
		conf:     *conf.New(),
		lang:     lang.New(lang.RUS),
		version:  "1.0.0",
	}

	// 2. Восстановление паники
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf(ac.T("ПАНИКА: %v\n"), r)
			// Выводим стек вызовов
			debug.PrintStack()
		}
	}()

	// 3. Создаем корневой контекст для graceful shutdown
	var rootCancel context.CancelFunc
	ac.rootCtx, rootCancel = context.WithCancel(context.Background())
	defer rootCancel()

	// 4. Инициализируем логгер приложения
	err := initLogger(&ac)
	if err != nil {
		fmt.Printf(ac.T("Ошибка создания логгера: %v\n"), err)
		os.Exit(1)
	}
	defer ac.logger.Close()

}
