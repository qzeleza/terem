package terem

import (
	"time"

	log "github.com/qzeleza/zlogger"
)

func initLogger(appConf *appConfig) error {

	// Создаем конфигурацию с настройками по умолчанию
	logConfig := log.NewConfig(appConf.logFile, appConf.sockFile)

	// Настраиваем дополнительные параметры
	if appConf.conf.LogMode == "true" {
		logConfig.Level = "debug"
	} else {
		logConfig.Level = "info"
	}
	logConfig.MaxFileSize = 50                // максимальный размер файла логов
	logConfig.BufferSize = 500                // размер буфера
	logConfig.FlushInterval = 2 * time.Second // интервал обновления буфера

	// Создаем логгер для основного приложения
	logger, err := log.New(logConfig, "MAIN")
	if err != nil {
		return err
	}
	defer logger.Close()

	// Возвращаем экземпляр логгера
	appConf.logger = *logger

	return nil
}
