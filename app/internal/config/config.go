package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qzeleza/terem/internal/utils"
	"gopkg.in/yaml.v3"
)

const (
	defaultConfigRoot      = "/opt/etc"
	defaultConfigDirectory = "terem"
	defaultConfigFile      = "config.yaml"
	configEnvVariable      = "TEREM_CONFIG"
	defaultLogFilePath     = "/tmp/terem.log"
)

// Config описывает настройки приложения. Та же структура сохраняется в YAML.
type Config struct {
	DebugMode bool   `yaml:"debugMode" json:"debugMode"` // Режим отладки
	LogFile   string `yaml:"logFile" json:"logFile"`     // Путь до файла логов
}

// Load загружает конфигурацию и гарантирует наличие файлов/директорий.
// Если путь недоступен, используетсяFallback в /tmp.
// explicitPath — явный путь до файла конфигурации.
func Load(explicitPath string) (*Config, string, error) {
	cfg := defaultConfig()

	path, err := resolveConfigPath(explicitPath)
	if err != nil {
		return nil, "", fmt.Errorf("определение пути конфигурации: %w", err)
	}

	path, err = ensureConfigFile(path, cfg)
	if err != nil {
		return nil, "", fmt.Errorf("создание конфигурационного файла: %w", err)
	}

	changedByMerge, err := mergeConfigFromFile(path, cfg)
	if err != nil {
		return nil, "", fmt.Errorf("чтение конфигурационного файла: %w", err)
	}

	previousLogPath := cfg.LogFile
	cfg.LogFile = ensureLogFilePath(cfg.LogFile)
	logPathAdjusted := cfg.LogFile != previousLogPath

	if changedByMerge || logPathAdjusted {
		_ = writeConfigFile(path, cfg)
	}

	return cfg, path, nil
}

// MustLoad игнорирует ошибки и возвращает конфигурацию по умолчанию в случае сбоя.
// explicitPath — явный путь до файла конфигурации.
func MustLoad(explicitPath string) *Config {
	cfg, _, err := Load(explicitPath)
	if err != nil {
		return defaultConfig()
	}
	return cfg
}

// SetDebugMode обновляет значение debugMode.
// v — новое значение debugMode.
func (c *Config) SetDebugMode(v bool) {
	if c == nil {
		return
	}
	c.DebugMode = v
}

// SetLogFile обновляет путь до файла логов.
// path — путь до файла логов.
func (c *Config) SetLogFile(path string) {
	if c == nil {
		return
	}
	c.LogFile = ensureLogFilePath(path)
}

// defaultConfig возвращает конфигурацию по умолчанию.
func defaultConfig() *Config {
	return &Config{
		DebugMode: utils.GetEnvBool("DEBUG", true),
		LogFile:   utils.GetEnv("TEREM_LOG_FILE", defaultLogFilePath),
	}
}

// resolveConfigPath возвращает путь до файла конфигурации.
// Если указан явный путь, возвращает его.
// Если переменная окружения TEREM_CONFIG установлена, возвращает её значение.
// Иначе возвращает путь по умолчанию.
// explicit — явный путь до файла конфигурации.
func resolveConfigPath(explicit string) (string, error) {
	switch {
	case explicit != "":
		return explicit, nil
	case os.Getenv(configEnvVariable) != "":
		return os.Getenv(configEnvVariable), nil
	}

	base := filepath.Join(defaultConfigRoot, defaultConfigDirectory)
	return filepath.Join(base, defaultConfigFile), nil
}

// ensureConfigFile создает конфигурационный файл, если он отсутствует.
// path — путь до файла конфигурации.
// cfg — конфигурация.
func ensureConfigFile(path string, cfg *Config) (string, error) {
	candidate := path

	if err := ensureDirectory(filepath.Dir(candidate)); err != nil {
		candidate = fallbackConfigPath(candidate)
		if err := ensureDirectory(filepath.Dir(candidate)); err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(candidate); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			candidate = fallbackConfigPath(candidate)
			if err := ensureDirectory(filepath.Dir(candidate)); err != nil {
				return "", err
			}
		}
		if err := writeConfigFile(candidate, cfg); err != nil {
			candidate = fallbackConfigPath(candidate)
			if err := writeConfigFile(candidate, cfg); err != nil {
				return "", err
			}
		}
	}

	return candidate, nil
}

// ensureDirectory создает директорию, если она отсутствует.
// dir — директория.
func ensureDirectory(dir string) error {
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

// fallbackConfigPath возвращает путь до файла конфигурации в директории /tmp.
// original — оригинальный путь до файла конфигурации.
func fallbackConfigPath(original string) string {
	base := filepath.Base(original)
	return filepath.Join(os.TempDir(), base)
}

// mergeConfigFromFile загружает конфигурацию из файла.
// path — путь до файла конфигурации.
// cfg — конфигурация.
func mergeConfigFromFile(path string, cfg *Config) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	var fileCfg Config
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return false, err
	}

	originalDebug := cfg.DebugMode
	originalLog := cfg.LogFile

	cfg.DebugMode = fileCfg.DebugMode
	if fileCfg.LogFile != "" {
		cfg.LogFile = fileCfg.LogFile
	}

	changed := originalDebug != cfg.DebugMode || originalLog != cfg.LogFile
	return changed, nil
}

// ensureLogFilePath создает директорию для файла логов, если она отсутствует.
// path — путь до файла логов.
func ensureLogFilePath(path string) string {
	candidate := path
	if candidate == "" {
		candidate = defaultLogFilePath
	}

	dir := filepath.Dir(candidate)
	if dir == "" || dir == "." {
		candidate = filepath.Join(os.TempDir(), filepath.Base(candidate))
		dir = filepath.Dir(candidate)
	}

	if err := ensureDirectory(dir); err != nil {
		candidate = filepath.Join(os.TempDir(), filepath.Base(candidate))
		_ = ensureDirectory(filepath.Dir(candidate))
	}

	if f, err := os.OpenFile(candidate, os.O_CREATE|os.O_APPEND, 0o644); err == nil {
		_ = f.Close()
		return candidate
	}

	fallback := filepath.Join(os.TempDir(), filepath.Base(candidate))
	_ = ensureDirectory(filepath.Dir(fallback))
	if f, err := os.OpenFile(fallback, os.O_CREATE|os.O_APPEND, 0o644); err == nil {
		_ = f.Close()
		return fallback
	}

	return candidate
}

// writeConfigFile записывает конфигурацию в файл.
// path — путь до файла конфигурации.
// cfg — конфигурация.
func writeConfigFile(path string, cfg *Config) error {
	if err := ensureDirectory(filepath.Dir(path)); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// MarshalJSON сериализует конфигурацию в JSON.
func (c *Config) MarshalJSON() ([]byte, error) {
	type alias Config
	return json.Marshal((*alias)(c))
}

// Save записывает конфигурацию в файл.
// path — путь до файла конфигурации.
func (c *Config) Save(path string) error {
	if c == nil {
		return errors.New("конфигурация не инициализирована")
	}
	if path == "" {
		return errors.New("путь к файлу не указан")
	}

	path, err := ensureConfigFile(path, c)
	if err != nil {
		return err
	}

	c.LogFile = ensureLogFilePath(c.LogFile)
	return writeConfigFile(path, c)
}
