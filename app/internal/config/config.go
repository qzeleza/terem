package config

import (
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
)

// Config описывает настройки приложения. Та же структура сериализуется в YAML-файл.
type Config struct {
	DebugMode bool   `yaml:"debugMode"`
	LogFile   string `yaml:"logFile"`
}

// Load возвращает конфигурацию (с учётом дефолтов) и путь до файла конфигурации.
// Если файл отсутствует — создаёт его с настройками по умолчанию в YAML-формате.
// explicitPath — явный путь до файла конфигурации.
func Load(explicitPath string) (*Config, string, error) {
	cfg := defaultConfig()

	// Определение пути конфигурации
	path, err := resolveConfigPath(explicitPath)
	if err != nil {
		return nil, "", fmt.Errorf("определение пути конфигурации: %w", err)
	}

	if err := ensureConfigFile(path, cfg); err != nil {
		return nil, "", fmt.Errorf("создание конфигурационного файла: %w", err)
	}

	if err := mergeConfigFromFile(path, cfg); err != nil {
		return nil, "", fmt.Errorf("чтение конфигурационного файла: %w", err)
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

// SetDebugMode обновляет значение debugMode и сразу отражает его в структуре.
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
	c.LogFile = path
}

// defaultConfig возвращает конфигурацию по умолчанию.
func defaultConfig() *Config {
	return &Config{
		DebugMode: utils.GetEnvBool("DEBUG", true),
		LogFile:   utils.GetEnv("TEREM_LOG_FILE", "/tmp/terem.log"),
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
func ensureConfigFile(path string, cfg *Config) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return writeConfigFile(path, cfg)
}

// mergeConfigFromFile загружает конфигурацию из файла.
// path — путь до файла конфигурации.
// cfg — конфигурация.
func mergeConfigFromFile(path string, cfg *Config) error {
	// Чтение файла конфигурации
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Разбор файла конфигурации
	var fileCfg Config
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return err
	}

	// Обновление конфигурации
	if fileCfg.LogFile != "" {
		cfg.LogFile = fileCfg.LogFile
	}
	cfg.DebugMode = fileCfg.DebugMode

	return nil
}

// writeConfigFile записывает конфигурацию в файл.
// path — путь до файла конфигурации.
// cfg — конфигурация.
func writeConfigFile(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Save перезаписывает конфигурационный файл.
// path — путь до файла конфигурации.
func (c *Config) Save(path string) error {
	if c == nil {
		return errors.New("конфигурация не инициализирована")
	}
	if path == "" {
		return errors.New("путь к файлу не указан")
	}
	return writeConfigFile(path, c)
}
