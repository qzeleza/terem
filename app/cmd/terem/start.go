package terem

import (
	"time"

	log "github.com/qzeleza/zlogger"
)

func initLogger(appConf *appConfig) error {

	// Создаем конфигурацию с настройками по умолчанию
	logConfig := log.NewConfig(appConf.logFile, appConf.sockFile)

	// Настраиваем дополнительные параметры
	logConfig.Level = appConf.conf.LogMode
	logConfig.MaxFileSize = 50 // 50 MB
	logConfig.BufferSize = 500
	logConfig.FlushInterval = 2 * time.Second

	// Создаем логгер с дополнительными сервисами
	logger, err := log.New(logConfig, "MAIN")
	if err != nil {
		return err
	}
	defer logger.Close()

	// Возвращаем экземпляр логгера
	appConf.logger = *logger

	return nil
}
