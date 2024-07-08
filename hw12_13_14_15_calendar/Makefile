BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# скомпилировать бинарный файл сервиса
build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

# собрать и запустить сервис с конфигом по умолчанию
run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

# запустить юнит-тесты
test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

# запустить golangci-lint 
lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

# применить миграции
.PHONY: migrate
migrate:
	goose -dir migrations postgres "postgres://otus:password@localhost:5432/calendar" up

# откатить миграции
.PHONY: migrate-down
migrate-down:
	goose -dir migrations postgres "postgres://otus:password@localhost:5432/calendar" down
