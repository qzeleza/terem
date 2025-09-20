# zlog Logger Module

`app/internal/zlog` реализует компактный логгер поверх [`zerolog`](https://github.com/rs/zerolog) с файловой ротацией (`lumberjack`) и единым текстовым форматом записей:

```
02-01-2006 15:04:05 [INFO] сообщение
```

Каждая запись содержит локальное время (ДД-ММ-ГГГГ ЧЧ:ММ:СС), уровень в квадратных скобках и текст сообщения. Ниже приведён обзор функций и примеры их применения.

## Создание и базовая настройка

### `New(filename string) *Logger`

Создаёт новый экземпляр логгера, привязанный к файлу `filename`. По умолчанию уровни форматируются как `[DEBUG]`, `[INFO]` и т.д., запись ведётся без цветов, а параметры ротации позже уточняются `AutoProfile`.

```go
logger := zlog.New("/tmp/terem.log")
logger.Info("Старт приложения") // 02-01-2006 15:04:05 [INFO] Старт приложения
```

### `(*Logger) AutoProfile()`

Автоматически подбирает параметры ротации (`MaxSize`, `MaxBackups`, `MaxAge`, `Compress`) на основании объёма RAM устройства (читается `/proc/meminfo`). Метод уже вызывается в `New`, но может быть вызван вручную повторно.

```go
logger.AutoProfile() // подстроит лимиты под текущую платформу
```

### Уровни логирования

Константы `DebugLevel`, `InfoLevel`, `WarnLevel`, `ErrorLevel`, `FatalLevel`, `PanicLevel`, `NoLevel`, `Disabled` соответствуют уровням `zerolog.Level`.

```go
logger.SetLevel(zlog.InfoLevel) // отключит DEBUG и ниже
```

## Управление ротацией

Все операции настраивают экземпляр `lumberjack.Logger` и влияют на формат файла: в выдачу самих логов изменения не вносятся, только регламентируют объём и срок хранения.

### `(*Logger) SetMaxSize(size int)`

Максимальный размер файла в мегабайтах, после чего запускается ротация.

```go
logger.SetMaxSize(20) // 20 MB на файл
```

### `(*Logger) SetMaxBackups(backups int)`

Количество резервных файлов, которые хранятся после ротации.

```go
logger.SetMaxBackups(5)
```

### `(*Logger) SetMaxAge(age int)`

Максимальный возраст резервных файлов в днях.

```go
logger.SetMaxAge(30) // хранить месяц
```

### `(*Logger) SetCompress(compress bool)`

Включает gzip-сжатие резервных файлов.

```go
logger.SetCompress(true)
```

### `(*Logger) Rotate() error`

Принудительно закрывает текущий файл и создаёт новый, возвращает ошибку, если операция не удалась.

```go
if err := logger.Rotate(); err != nil {
    logger.Error("Не удалось ротировать лог: %v", err)
}
```

### `(*Logger) Close() error`

Закрывает файл, удаляет логгер из глобального реестра, чтобы он больше не участвовал в обработке SIGHUP. Важно вызывать перед завершением программы, если логгер более не нужен.

```go
if err := logger.Close(); err != nil {
    fmt.Printf("ошибка закрытия: %v\n", err)
}
```

## Логирование событий

Все методы возвращают `error` для совместимости с прежним API, но в текущей реализации всегда возвращают `nil`. Формат вывода одинаков: `02-01-2006 15:04:05 [LEVEL] сообщение`. Форматирование реализовано через `fmt.Sprintf`, если первый аргумент — строка формата.

### `(*Logger) Debug(args ...interface{}) error`

```go
logger.Debug("raw payload: %x", payload)
// => 02-01-2006 15:04:05 [DEBUG] raw payload: DEADBEEF
```

### `(*Logger) Info(args ...interface{}) error`

```go
logger.Info("Запущен профиль %s", profileName)
// => 02-01-2006 15:04:05 [INFO] Запущен профиль default
```

### `(*Logger) Warn(args ...interface{}) error`

```go
logger.Warn("Низкий уровень памяти: %d MB", freeMem)
// => 02-01-2006 15:04:05 [WARN] Низкий уровень памяти: 42 MB
```

### `(*Logger) Error(args ...interface{}) error`

```go
logger.Error("Ошибка подключения: %v", err)
// => 02-01-2006 15:04:05 [ERROR] Ошибка подключения: dial timeout
```

### `(*Logger) Fatal(args ...interface{}) error`

Записывает сообщение и завершает программу через `os.Exit(1)`.

```go
logger.Fatal("Критическая ошибка конфигурации")
// => 02-01-2006 15:04:05 [FATAL] Критическая ошибка конфигурации (после записи приложение завершится)
```

### `(*Logger) Panic(args ...interface{}) error`

Пишет запись и вызывает `panic`.

```go
logger.Panic("Нарушена целостность системных файлов")
// => 02-01-2006 15:04:05 [PANIC] Нарушена целостность системных файлов
```

## Работа с сигналами

### `EnableSIGHUP()`

Однократно настраивает обработчик `SIGHUP`. При получении сигнала выполняется `doRotate()` — ротация всех зарегистрированных логгеров. Сообщения о результате попадают в лог в формате `[INFO]` или `[ERROR]`. Если ни одного логгера не зарегистрировано, сообщения выводятся в `stderr`.

```go
func main() {
    zlog.EnableSIGHUP()
    logger := zlog.New("/var/log/terem.log")
    defer logger.Close()
    select {} // приложение продолжает работать, ожидая SIGHUP
}
```

## Внутренние вспомогательные функции

Хотя они не экспортируются, их важно понимать при расширении пакета.

- `registerLogger(*Logger)` / `unregisterLogger(*Logger)` — добавляют/удаляют логгеры из реестра, чтобы групповые операции (`doRotate`) знали о них.
- `snapshotLoggers() []*Logger` — потокобезопасная копия реестра для итерирования без блокировки на запись.
- `primaryLogger() *Logger` — первый доступный логгер, используется для сообщений об ошибках/успехе ротации по SIGHUP.
- `doRotate() error` — запускает `Rotate()` для всех логгеров, собирает ошибки через `errors.Join` и логирует результат через `logSIGHUPError`/`logSIGHUPSuccess`.
- `logSIGHUPError(error)` и `logSIGHUPSuccess()` — оформляют сообщения в формате `02-01-2006 15:04:05 [ERROR] ...` или `[..., INFO] ...`, либо печатают в stderr, если логгеров нет.
- `formatArgs(args ...interface{}) string` — общая функция форматирования, поддерживает `fmt`-совместимые шаблоны.
- `getTotalMemory() uint64` — читает `/proc/meminfo` и возвращает объём RAM; используется в `AutoProfile` для выбора параметров ротации.

## Полезные комбинации

```go
func setupLogger(path string, debug bool) (*zlog.Logger, error) {
    logger := zlog.New(path)
    if debug {
        logger.SetLevel(zlog.DebugLevel)
    } else {
        logger.SetLevel(zlog.InfoLevel)
    }
    logger.SetMaxSize(25)
    logger.SetMaxBackups(7)
    logger.SetCompress(true)
    zlog.EnableSIGHUP()
    logger.Info("Логгер готов к работе")
    return logger, nil
}
```

Этот пример создаёт логгер с параметрами ротации, текстовым выводом и обработкой SIGHUP. Каждая запись продолжает следовать формату `02-01-2006 15:04:05 [LEVEL] сообщение`.
