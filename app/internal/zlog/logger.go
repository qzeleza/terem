package zlog

// Пакет logger предоставляет лёгкий, высокопроизводительный логгер на основе zerolog с ротацией файлов (lumberjack).
// Формат записей: 02-01-2006 15:04:05 [LEVEL] сообщение. Подходит для embedded-устройств с минимальным overhead.

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/qzeleza/terem/internal/i18n"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

const consoleTimeFormat = "02-01-2006 15:04:05"

// Logger — структура логгера с настройками ротации.
// Поле zlog — внутренний zerolog.Logger; rot — ротационный writer.
type Logger struct {
	zlog zerolog.Logger
	rot  *lumberjack.Logger
}

var (
	registry   = newLoggerRegistry()
	sighupOnce sync.Once
)

// Уровни логирования (константы пакета, равны zerolog.Level для удобства использования как log.InfoLevel).
const (
	DebugLevel zerolog.Level = zerolog.DebugLevel
	InfoLevel  zerolog.Level = zerolog.InfoLevel
	WarnLevel  zerolog.Level = zerolog.WarnLevel
	ErrorLevel zerolog.Level = zerolog.ErrorLevel
	FatalLevel zerolog.Level = zerolog.FatalLevel
	PanicLevel zerolog.Level = zerolog.PanicLevel
	NoLevel    zerolog.Level = zerolog.NoLevel
	Disabled   zerolog.Level = zerolog.Disabled
)

// New создаёт новый экземпляр логгера и регистрирует его в глобальном реестре.
// Параметры ротации подстраиваются автоматически через AutoProfile.
func New(filename string) *Logger {
	rotator := newRotator(filename)
	writer := newConsoleWriter(rotator)
	z := newZerolog(writer)

	logger := &Logger{zlog: z, rot: rotator}
	registry.Add(logger)
	logger.AutoProfile()

	return logger
}

// AutoProfile устанавливает параметры ротации на основе объёма памяти устройства.
func (l *Logger) AutoProfile() {
	if l == nil || l.rot == nil {
		return
	}

	mem := getTotalMemory()

	switch {
	case mem < 65*1024*1024:
		l.rot.MaxSize = 5
		l.rot.MaxBackups = 1
		l.rot.MaxAge = 3
		l.rot.Compress = false
	case mem < 128*1024*1024:
		l.rot.MaxSize = 10
		l.rot.MaxBackups = 2
		l.rot.MaxAge = 7
		l.rot.Compress = false
	default:
		l.rot.MaxSize = 50
		l.rot.MaxBackups = 3
		l.rot.MaxAge = 14
		l.rot.Compress = true
	}
}

// SetLevel устанавливает глобальный уровень логирования (использует константы пакета).
func (l *Logger) SetLevel(level zerolog.Level) {
	if l == nil {
		return
	}
	l.zlog = l.zlog.Level(level)
}

// SetMaxSize устанавливает максимальный размер файла перед ротацией (в MB).
func (l *Logger) SetMaxSize(size int) {
	if l == nil || l.rot == nil {
		return
	}
	l.rot.MaxSize = size
}

// SetMaxBackups устанавливает максимальное количество резервных файлов.
func (l *Logger) SetMaxBackups(backups int) {
	if l == nil || l.rot == nil {
		return
	}
	l.rot.MaxBackups = backups
}

// SetMaxAge устанавливает максимальный возраст резервных файлов (в днях).
func (l *Logger) SetMaxAge(age int) {
	if l == nil || l.rot == nil {
		return
	}
	l.rot.MaxAge = age
}

// SetCompress включает/выключает сжатие резервных файлов (gzip).
func (l *Logger) SetCompress(compress bool) {
	if l == nil || l.rot == nil {
		return
	}
	l.rot.Compress = compress
}

// Rotate принудительно выполняет ротацию текущего файла лога.
func (l *Logger) Rotate() error {
	if l == nil || l.rot == nil {
		return nil
	}
	return l.rot.Rotate()
}

// Close закрывает writer и исключает логгер из глобального реестра.
func (l *Logger) Close() error {
	if l == nil {
		return nil
	}
	registry.Remove(l)
	if l.rot == nil {
		return nil
	}
	return l.rot.Close()
}

// EnableSIGHUP включает обработчик SIGHUP для всех зарегистрированных логгеров.
func EnableSIGHUP() {
	sighupOnce.Do(func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP)
		go func() {
			for range c {
				if err := registry.RotateAll(); err != nil {
					logSIGHUPError(err)
				} else {
					logSIGHUPSuccess()
				}
			}
		}()
	})
}

// Debug логирует на уровне Debug. Возвращает nil для совместимости с прежним API.
func (l *Logger) Debug(args ...interface{}) error {
	if l == nil {
		return nil
	}
	l.zlog.Debug().Msg(formatArgs(args...))
	return nil
}

// Info логирует на уровне Info. Возвращает nil для совместимости с прежним API.
func (l *Logger) Info(args ...interface{}) error {
	if l == nil {
		return nil
	}
	l.zlog.Info().Msg(formatArgs(args...))
	return nil
}

// Warn логирует на уровне Warn. Возвращает nil для совместимости с прежним API.
func (l *Logger) Warn(args ...interface{}) error {
	if l == nil {
		return nil
	}
	l.zlog.Warn().Msg(formatArgs(args...))
	return nil
}

// Error логирует на уровне Error. Возвращает nil для совместимости с прежним API.
func (l *Logger) Error(args ...interface{}) error {
	if l == nil {
		return nil
	}
	l.zlog.Error().Msg(formatArgs(args...))
	return nil
}

// Fatal логирует на уровне Fatal и завершает приложение.
func (l *Logger) Fatal(args ...interface{}) error {
	if l == nil {
		return nil
	}
	l.zlog.Fatal().Msg(formatArgs(args...))
	return nil
}

// Panic логирует на уровне Panic и вызывает панику.
func (l *Logger) Panic(args ...interface{}) error {
	if l == nil {
		return nil
	}
	l.zlog.Panic().Msg(formatArgs(args...))
	return nil
}

func newRotator(filename string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     14,
		Compress:   true,
		LocalTime:  true,
	}
}

func newConsoleWriter(out io.Writer) zerolog.ConsoleWriter {
	writer := zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: consoleTimeFormat,
		NoColor:    true,
	}
	writer.FormatLevel = func(i interface{}) string {
		return fmt.Sprintf("[%s] ", strings.ToUpper(fmt.Sprint(i)))
	}
	writer.FormatMessage = func(i interface{}) string {
		return fmt.Sprint(i)
	}
	writer.FormatFieldName = func(i interface{}) string { return "" }
	writer.FormatFieldValue = func(i interface{}) string { return "" }

	return writer
}

func newZerolog(writer zerolog.ConsoleWriter) zerolog.Logger {
	return zerolog.New(writer).With().Timestamp().Logger().Level(DebugLevel)
}

func logSIGHUPError(err error) {
	if log := registry.Primary(); log != nil {
		log.zlog.Error().Err(err).Msg(i18n.T("logger.error.force_rotate"))
		return
	}
	fmt.Fprintf(os.Stderr, i18n.T("logger.error.force_rotate_stderr")+"\n", err)
}

func logSIGHUPSuccess() {
	if log := registry.Primary(); log != nil {
		log.zlog.Info().Msg(i18n.T("logger.info.force_rotate"))
		return
	}
	fmt.Fprintln(os.Stderr, i18n.T("logger.info.force_rotate_stderr"))
}

type loggerRegistry struct {
	mu      sync.RWMutex
	loggers map[*Logger]struct{}
}

func newLoggerRegistry() *loggerRegistry {
	return &loggerRegistry{loggers: make(map[*Logger]struct{})}
}

func (r *loggerRegistry) Add(l *Logger) {
	if l == nil {
		return
	}
	r.mu.Lock()
	r.loggers[l] = struct{}{}
	r.mu.Unlock()
}

func (r *loggerRegistry) Remove(l *Logger) {
	if l == nil {
		return
	}
	r.mu.Lock()
	delete(r.loggers, l)
	r.mu.Unlock()
}

func (r *loggerRegistry) Snapshot() []*Logger {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.loggers) == 0 {
		return nil
	}

	items := make([]*Logger, 0, len(r.loggers))
	for logger := range r.loggers {
		items = append(items, logger)
	}
	return items
}

func (r *loggerRegistry) Primary() *Logger {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for logger := range r.loggers {
		return logger
	}
	return nil
}

func (r *loggerRegistry) RotateAll() error {
	snapshot := r.Snapshot()
	if len(snapshot) == 0 {
		return nil
	}

	var errs []error
	for _, logger := range snapshot {
		if err := logger.Rotate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
