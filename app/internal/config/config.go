// Package config содержит конфигурацию приложения
package config

import (
	"os"
)

// Config представляет конфигурацию приложения
type Config struct {
	// Архитектура устройства
	Arch string
	// Режим работы (development, production)
	DevMode string
	// Режим логирования
	DebugMode string
}

// New создаёт новую конфигурацию с значениями по умолчанию
func New() *Config {
	return &Config{
		Arch:      getEnv("ARCH", "arm"),
		DevMode:   getEnv("MODE", "develop"),
		DebugMode: getEnv("DEBUG", "true"),
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
