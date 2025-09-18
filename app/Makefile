# Переменные
APP_NAME=terem
MAIN_PATH=cmd/$(APP_NAME)/main.go
BUILD_DIR=../builder

# Команды по умолчанию
.DEFAULT_GOAL := help

## help: показать это сообщение
help:
	@echo "Доступные команды:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run: запустить приложение
run:
	go run $(MAIN_PATH)

## build: собрать приложение
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

## test: запустить тесты
test:
	go test -v ./...

## test-coverage: запустить тесты с покрытием
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## clean: очистить сборку
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

## deps: установить зависимости
deps:
	go clean -modcache
	go mod download
	go mod tidy

## fmt: форматировать код
fmt:
	go fmt ./...

## lint: проверить код линтером
lint:
	golangci-lint run

## cache: очистить кэш
cache:
	go clean -cache
	go clean -testcache
	go clean -modcache

## install: установить приложение
install: build
	cp $(BUILD_DIR)/$(APP_NAME) $(GOPATH)/bin/

.PHONY: help run build test test-coverage clean deps fmt lint cache install