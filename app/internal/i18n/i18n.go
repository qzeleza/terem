package i18n

import (
	"bufio"
	"fmt"
	"strings"
	"sync"

	"github.com/qzeleza/terem/internal/lang"
)

const fallbackLanguage = "ru"

var (
	loadOnce        sync.Once
	loadErr         error
	dictionaries    = map[string]map[string]string{}
	currentLanguage = fallbackLanguage
	mu              sync.RWMutex
)

// Available возвращает список поддерживаемых языков.
func Available() []string {
	ensureLoaded()

	mu.RLock()
	defer mu.RUnlock()

	codes := make([]string, 0, len(dictionaries))
	for code := range dictionaries {
		codes = append(codes, code)
	}
	return codes
}

// SetLanguage выбирает текущий язык. При неизвестном языке возвращает ошибку.
func SetLanguage(lang string) error {
	ensureLoaded()

	mu.Lock()
	defer mu.Unlock()

	if _, ok := dictionaries[lang]; !ok {
		return fmt.Errorf("language %s is not available", lang)
	}

	currentLanguage = lang
	return nil
}

// Language возвращает текущий язык.
func Language() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentLanguage
}

// T возвращает строку для ключа key. При отсутствии ключа возвращает его же.
func T(key string, args ...any) string {
	ensureLoaded()

	mu.RLock()
	lang := currentLanguage
	langDict := dictionaries[lang]
	fallbackDict := dictionaries[fallbackLanguage]
	mu.RUnlock()

	if value, ok := langDict[key]; ok {
		if len(args) > 0 {
			return fmt.Sprintf(value, args...)
		}
		return value
	}

	if value, ok := fallbackDict[key]; ok {
		if len(args) > 0 {
			return fmt.Sprintf(value, args...)
		}
		return value
	}

	return key
}

// ensureLoaded загружает словари единожды.
func ensureLoaded() {
	loadOnce.Do(func() {
		loadErr = loadAll()
	})
}

// Error возвращает ошибку, возникшую при загрузке словарей.
func Error() error {
	ensureLoaded()
	return loadErr
}

func loadAll() error {
	entries, err := lang.Files.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read language directory: %w", err)
	}

	tmp := make(map[string]map[string]string, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		code := name
		if strings.HasSuffix(name, ".txt") {
			code = strings.TrimSuffix(name, ".txt")
		} else if strings.Contains(name, ".") {
			continue
		}
		if code == "" {
			continue
		}
		data, err := lang.Files.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}

		dict, err := parseDictionary(string(data))
		if err != nil {
			return fmt.Errorf("parse %s: %w", name, err)
		}
		tmp[code] = dict
	}

	mu.Lock()
	dictionaries = tmp
	if _, ok := dictionaries[currentLanguage]; !ok {
		currentLanguage = fallbackLanguage
	}
	mu.Unlock()

	if _, ok := dictionaries[fallbackLanguage]; !ok {
		return fmt.Errorf("fallback language %s is missing", fallbackLanguage)
	}

	return nil
}

func parseDictionary(data string) (map[string]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	result := make(map[string]string)
	lineNumber := 0

	for scanner.Scan() {
		rawLine := scanner.Text()
		lineNumber++

		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		key, value, ok := splitKeyValue(line)
		if !ok {
			return nil, fmt.Errorf("line %d: missing key or value", lineNumber)
		}

		key = strings.TrimSpace(unescape(key))
		value = strings.TrimSpace(unescape(value))
		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNumber)
		}
		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func splitKeyValue(line string) (string, string, bool) {
	escaped := false
	for i, r := range line {
		if r == '=' && !escaped {
			return line[:i], line[i+1:], true
		}
		if r == '\\' && !escaped {
			escaped = true
			continue
		}
		escaped = false
	}
	return "", "", false
}

func unescape(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))

	escaped := false
	for i := 0; i < len(input); i++ {
		ch := input[i]
		if escaped {
			switch ch {
			case 'n':
				builder.WriteByte('\n')
			case 't':
				builder.WriteByte('\t')
			case 'r':
				builder.WriteByte('\r')
			case '=':
				builder.WriteByte('=')
			case '\\':
				builder.WriteByte('\\')
			default:
				builder.WriteByte(ch)
			}
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		builder.WriteByte(ch)
	}

	if escaped {
		builder.WriteByte('\\')
	}

	return builder.String()
}
