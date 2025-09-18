// Package lang предоставляет функциональность для локализации приложения.
package lang

import (
	"fmt"
	"strings"

	"github.com/qzeleza/terem/internal/lang/translation"
)

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

var AvailableLangs = map[Language]string{
	RUS: "Русский",
	ENG: "English",
	BEL: "Беларуская",
	KAZ: "Қазақша",
}

// Translator - это структура, которая хранит выбранный язык
// и предоставляет метод для перевода строк.
type Translator struct {
	language Language
}

// Убеждаемся, что Translator реализует TranslatorInterface
var _ TranslatorInterface = (*Translator)(nil)

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
		dictionary = translation.Eng
	case BEL:
		dictionary = translation.Bel
	case KAZ:
		dictionary = translation.Kaz
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

// ParseLanguage преобразует строку языка в enum
func ParseLanguage(langStr string) (Language, error) {
	langStr = strings.ToLower(strings.TrimSpace(langStr))

	switch langStr {
	case "ru", "rus", "русский", "российский":
		return RUS, nil
	case "en", "eng", "english", "английский":
		return ENG, nil
	case "be", "bel", "беларуская", "белорусский":
		return BEL, nil
	case "kz", "kaz", "қазақша", "казахский":
		return KAZ, nil
	default:
		return RUS, fmt.Errorf("неподдерживаемый язык: %s", langStr)
	}
}
