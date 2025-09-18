package lang

import (
	"sync"
)

// OptimizedTranslator - оптимизированная версия переводчика для роутеров
// Использует ленивую загрузку и кэширование только нужных переводов
type OptimizedTranslator struct {
	language    Language
	cache       map[string]string
	cacheMutex  sync.RWMutex
	maxCache    int // максимальный размер кэша
}

// Убеждаемся, что OptimizedTranslator реализует TranslatorInterface
var _ TranslatorInterface = (*OptimizedTranslator)(nil)

// NewOptimized создает оптимизированный переводчик с ограниченным кэшем
func NewOptimized(lang Language, maxCacheSize int) *OptimizedTranslator {
	if maxCacheSize <= 0 {
		maxCacheSize = 50 // разумный лимит для роутеров
	}

	return &OptimizedTranslator{
		language: lang,
		cache:    make(map[string]string, maxCacheSize),
		maxCache: maxCacheSize,
	}
}

// T выполняет перевод с ленивой загрузкой и кэшированием
func (t *OptimizedTranslator) T(rusString string) string {
	// Если русский язык, возвращаем как есть
	if t.language == RUS {
		return rusString
	}

	// Проверяем кэш
	t.cacheMutex.RLock()
	if cached, exists := t.cache[rusString]; exists {
		t.cacheMutex.RUnlock()
		return cached
	}
	t.cacheMutex.RUnlock()

	// Получаем перевод и кэшируем
	translation := t.getTranslation(rusString)

	t.cacheMutex.Lock()
	defer t.cacheMutex.Unlock()

	// Проверяем размер кэша и очищаем при необходимости
	if len(t.cache) >= t.maxCache {
		// Простая стратегия - очищаем весь кэш (для роутеров)
		// В более сложных случаях можно использовать LRU
		t.cache = make(map[string]string, t.maxCache)
	}

	t.cache[rusString] = translation
	return translation
}

// getTranslation получает перевод без кэширования
func (t *OptimizedTranslator) getTranslation(rusString string) string {
	switch t.language {
	case ENG:
		return getEnglishTranslation(rusString)
	case BEL:
		return getBelarusianTranslation(rusString)
	case KAZ:
		return getKazakhTranslation(rusString)
	default:
		return rusString
	}
}

// ClearCache очищает кэш переводов (полезно для освобождения памяти)
func (t *OptimizedTranslator) ClearCache() {
	t.cacheMutex.Lock()
	defer t.cacheMutex.Unlock()
	t.cache = make(map[string]string, t.maxCache)
}

// GetCacheSize возвращает текущий размер кэша
func (t *OptimizedTranslator) GetCacheSize() int {
	t.cacheMutex.RLock()
	defer t.cacheMutex.RUnlock()
	return len(t.cache)
}