package lang

import (
	"os"
	"strconv"
)

// TranslatorFactory фабрика для создания переводчиков
type TranslatorFactory struct{}

// CreateTranslator создает оптимальный переводчик в зависимости от условий
func (tf *TranslatorFactory) CreateTranslator(lang Language) TranslatorInterface {
	// Определяем среду выполнения
	isLowMemory := tf.isLowMemoryEnvironment()

	if isLowMemory {
		// Для роутеров и встроенных систем - используем оптимизированную версию
		cacheSize := tf.getOptimalCacheSize()
		return NewOptimized(lang, cacheSize)
	}

	// Для обычных систем - используем стандартную версию
	return New(lang)
}

// TranslatorInterface общий интерфейс для всех переводчиков
type TranslatorInterface interface {
	T(rusString string) string
}

// isLowMemoryEnvironment определяет, работаем ли мы в среде с ограниченной памятью
func (tf *TranslatorFactory) isLowMemoryEnvironment() bool {
	// Проверяем переменные окружения
	if arch := os.Getenv("ARCH"); arch == "arm" || arch == "mips" {
		return true
	}

	// Проверяем доступную память (упрощенно)
	if memLimit := os.Getenv("MEMORY_LIMIT_MB"); memLimit != "" {
		if limit, err := strconv.Atoi(memLimit); err == nil && limit < 128 {
			return true
		}
	}

	// Проверяем режим разработки
	if devMode := os.Getenv("DEV_MODE"); devMode == "embedded" || devMode == "router" {
		return true
	}

	return false
}

// getOptimalCacheSize возвращает оптимальный размер кэша для текущей среды
func (tf *TranslatorFactory) getOptimalCacheSize() int {
	// Для роутеров - минимальный кэш
	if arch := os.Getenv("ARCH"); arch == "arm" || arch == "mips" {
		return 20
	}

	// Для обычных систем - больший кэш
	return 50
}

// NewFactory создает новую фабрику переводчиков
func NewFactory() *TranslatorFactory {
	return &TranslatorFactory{}
}