FROM golang:1.22-alpine as build

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar_integration_tests"
LABEL MAINTAINERS="student@otus.ru"

# Для использования временных зон
RUN apk add --no-cache tzdata

ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ENV CONFIG_FILE /etc/integration_tests/config.yaml

CMD go test -v -timeout 2m -count=1 --tags=integration /go/src/integration_tests/...