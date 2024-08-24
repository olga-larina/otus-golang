BIN_CALENDAR := "./bin/calendar"
DOCKER_IMG_CALENDAR="calendar:develop"

BIN_SCHEDULER := "./bin/calendar_scheduler"
DOCKER_IMG_SCHEDULER="calendar_scheduler:develop"

BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG_SENDER="calendar_sender:develop"

DOCKER_IMG_MIGRATOR="calendar_migrator:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# скомпилировать бинарные файлы сервиса
.PHONY: build_calendar
build_calendar:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar

.PHONY: build_scheduler
build_scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

.PHONY: build_sender
build_sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

.PHONY: build
build: build_calendar build_scheduler build_sender

# собрать и запустить сервисы с конфигами по умолчанию
.PHONY: run_calendar
run_calendar: build_calendar
	$(BIN_CALENDAR) -config ./configs/calendar_config.yaml

.PHONY: run_scheduler
run_scheduler: build_scheduler
	$(BIN_SCHEDULER) -config ./configs/scheduler_config.yaml
	
.PHONY: run_sender
run_sender: build_sender
	$(BIN_SENDER) -config ./configs/sender_config.yaml

# собрать образ миграций
.PHONY: build-img-migrator
build-img-migrator:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_MIGRATOR) \
		-f build/migrator/Dockerfile .

# собрать все образы
.PHONY: build-img
build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_CALENDAR) \
		-f build/calendar/Dockerfile .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SCHEDULER) \
		-f build/scheduler/Dockerfile .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SENDER) \
		-f build/sender/Dockerfile .

# запустить все образы
.PHONY: run-img
run-img: build-img
	docker run $(DOCKER_IMG_CALENDAR) -d
	docker run $(DOCKER_IMG_SCHEDULER) -d
	docker run $(DOCKER_IMG_SENDER) -d

# посмотреть версию приложения календарь
.PHONY: version
version: build_calendar
	$(BIN_CALENDAR) version

# запустить юнит-тесты
.PHONY: test
test:
	go test -race -count=1 ./internal/...

# линтер golangci-lint 
.PHONY: install-lint-deps
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

.PHONY: lint
lint: install-lint-deps
	golangci-lint run ./... --fix

# применить миграции
.PHONY: migrate
migrate:
	goose -dir migrations postgres "postgres://otus:password@localhost:5432/calendar" up

# откатить миграции
.PHONY: migrate-down
migrate-down:
	goose -dir migrations postgres "postgres://otus:password@localhost:5432/calendar" down

# сгенерировать прото
# необходимо сначала выполнить
# - go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# - go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# - создать пакет для кода grpc-сервера с go:generate
.PHONY: generate
generate: 
	go generate ./...

# поднять окружение
.PHONY: up-infra
up-infra:
	docker compose --env-file deployments/.env -f deployments/docker-compose-infra.yaml up -d

# потушить окружение
.PHONY: down-infra
down-infra:
	docker compose --env-file deployments/.env -f deployments/docker-compose-infra.yaml down

# поднять сервисы и окружение
.PHONY: up
up:
	docker compose --env-file deployments/.env -f deployments/docker-compose-infra.yaml -f deployments/docker-compose.yaml up -d --build

# потушить сервисы и окружение
.PHONY: down
down:
	docker compose --env-file deployments/.env -f deployments/docker-compose-infra.yaml -f deployments/docker-compose.yaml down

# запустить интеграционные тесты
.PHONY: integration-tests
integration-tests:
	bash integration_tests/integration_tests.sh; \
    e=$$?; \
    exit $$e
