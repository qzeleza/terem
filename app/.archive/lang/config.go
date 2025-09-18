package lang

import (
	"os"
	"strconv"
)

// Config конфигурация системы локализации
type Config struct {
	DefaultLanguage Language
	CacheSize       int
	EnableCache     bool
	CompactMode     bool // Использовать только компактные переводы
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		DefaultLanguage: RUS,
		CacheSize:       50,
		EnableCache:     true,
		CompactMode:     false,
	}
}

// RouterConfig возвращает оптимизированную конфигурацию для роутеров
func RouterConfig() *Config {
	return &Config{
		DefaultLanguage: RUS,
		CacheSize:       15, // Минимальный кэш
		EnableCache:     true,
		CompactMode:     true, // Только основные переводы
	}
}

// LoadConfigFromEnv загружает конфигурацию из переменных окружения
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	// Определяем архитектуру и режим работы
	if arch := os.Getenv("ARCH"); arch == "arm" || arch == "mips" {
		config = RouterConfig()
	}

	// Переопределяем из переменных окружения
	if langStr := os.Getenv("TEREM_LANG"); langStr != "" {
		config.DefaultLanguage, _ = ParseLanguage(langStr)
	}

	if cacheStr := os.Getenv("TEREM_CACHE_SIZE"); cacheStr != "" {
		if size, err := strconv.Atoi(cacheStr); err == nil && size > 0 {
			config.CacheSize = size
		}
	}

	if compactStr := os.Getenv("TEREM_COMPACT_MODE"); compactStr == "true" {
		config.CompactMode = true
	}

	return config
}

func SetLanguage(lang int) {
	DefaultConfig().DefaultLanguage = Language(lang)
}

// CreateTranslatorFromConfig создает переводчик на основе конфигурации
func CreateTranslatorFromConfig(config *Config) TranslatorInterface {
	if config.CompactMode || config.CacheSize < 30 {
		return NewOptimized(config.DefaultLanguage, config.CacheSize)
	}
	return New(config.DefaultLanguage)
}
