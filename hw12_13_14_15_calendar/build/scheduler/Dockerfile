# Собираем в гошке
FROM golang:1.22 AS build

ENV BIN_FILE /opt/scheduler/calendar_scheduler
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
RUN go mod tidy

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/scheduler/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar_scheduler"
LABEL MAINTAINERS="student@otus.ru"

# Для использования временных зон
RUN apk add --no-cache tzdata

ENV BIN_FILE "/opt/scheduler/calendar_scheduler"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/scheduler/config.yaml

CMD ${BIN_FILE} -config ${CONFIG_FILE}
