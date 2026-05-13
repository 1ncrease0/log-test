# Парсер логов

Подробные требования - в [TASK.md](TASK.md).

## Запуск

```bash
docker compose up --build -d
```

```bash
make up
```

Пример параметров в [.env.example](.env.example).


```bash
docker compose down
```

```bash
make down
make clean
```


## Тесты

```bash
make test
```

## Postman

Простые тесты [Parser.postman_collection.json](Parser.postman_collection.json) в Postman (после `docker compose up`). В коллекции есть тесты на ответы API.

## Примеры curl

Запуск парсинга (путь - относительно `data/`, внутри архива должны быть `ibdiagnet2.db_csv` и `ibdiagnet2.sharp_an_info`):

```bash
curl -sS -X POST http://localhost:8080/api/v1/parse/ \
  -H 'Content-Type: application/json' \
  -d '{"path":"log.zip"}'
```

Топология:

```bash
curl -sS http://localhost:8080/api/v1/topology/1
```

Информация лога:

```bash
curl -sS http://localhost:8080/api/v1/log/1
```

Узел и порты:

```bash
curl -sS http://localhost:8080/api/v1/node/1
curl -sS http://localhost:8080/api/v1/port/1
```

