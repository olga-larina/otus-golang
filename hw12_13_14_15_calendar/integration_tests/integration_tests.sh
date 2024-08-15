#!/bin/bash

# запускаем окружение и тесты с соответствующими переменными окружения
docker-compose --env-file deployments/.env-tests -f deployments/docker-compose.yaml -f deployments/docker-compose-tests.yaml up -d --build

# дожидаемся код ответа
exit_code=$(docker wait calendar_integration_tests) || exit_code=$?

# выводим логи тестов
docker logs calendar_integration_tests

# останавливаем с удалением volume
docker-compose --env-file deployments/.env-tests -f deployments/docker-compose.yaml -f deployments/docker-compose-tests.yaml down -v

echo "Exit code: $exit_code"

exit $exit_code