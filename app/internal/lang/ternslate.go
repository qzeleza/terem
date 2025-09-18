// Package lang предоставляет функциональность для локализации приложения.
package lang

// Language определяет тип для идентификаторов языка.
// Мы используем такой подход для создания подобия enum.
type Language int

// Определяем константы для каждого языка. iota автоматически присваивает
// последовательные целочисленные значения (0, 1, 2, ...).
const (
	RUS Language = iota // Русский (по умолчанию)
	ENG                 // Английский
	BEL                 // Белорусский
	KAZ                 // Казахский
)

// eng - словарь для английского языка.
// Ключ - фраза на русском, значение - перевод на английский.
var eng = map[string]string{
	"ПАНИКА": "PANIC",
	"Ошибка создания логгера": "Logger creation error",
}

// bel - словарь для белорусского языка.
var bel = map[string]string{
	"ПАНИКА": "ПАНІКА",
	"Ошибка создания логгера": "Памылка стварэння логера",
}

// kaz - словарь для казахского языка.
var kaz = map[string]string{
	"ПАНИКА": "ДҮРБЕЛЕҢ",
	"Ошибка создания логгера": "Логгерді құру қатесі",
}

// Translator - это структура, которая хранит выбранный язык
// и предоставляет метод для перевода строк.
type Translator struct {
	language Language
}

// New создает новый экземпляр Translator для указанного языка.
//
// @param lang Язык, который будет использоваться для перевода.
// @return *Translator Указатель на новый экземпляр Translator.
func New(lang Language) *Translator {
	return &Translator{language: lang}
}

// T выполняет перевод строки на язык, заданный при создании Translator.
// Если перевод для строки не найден, возвращается исходная строка.
//
// @param rusString Строка на русском языке для перевода.
// @return string Переведенная строка или исходная строка, если перевод не найден.
func (t *Translator) T(rusString string) string {
	// Если язык русский, просто возвращаем исходную строку.
	if t.language == RUS {
		return rusString
	}

	var dictionary map[string]string

	// Выбираем нужный словарь в зависимости от языка.
	switch t.language {
	case ENG:
		dictionary = eng
	case BEL:
		dictionary = bel
	case KAZ:
		dictionary = kaz
	default:
		// Если язык не поддерживается, возвращаем исходную строку.
		return rusString
	}

	// Ищем перевод в выбранном словаре.
	if translation, ok := dictionary[rusString]; ok {
		// Если нашли, возвращаем его.
		return translation
	}

	// Если не нашли, возвращаем исходную строку.
	return rusString
}
