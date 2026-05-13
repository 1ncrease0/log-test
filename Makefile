.PHONY: lint up down clean mocks test

MOCKERY ?= go run github.com/vektra/mockery/v2@v2.53.2

lint:
	golangci-lint run ./...

up: down
	docker compose up --build -d

down:
	docker compose down

clean:
	docker compose down -v

mocks:
	$(MOCKERY) --config .mockery.yaml

test:
	go test -race -coverprofile=cover.out \
		$(shell go list ./... | grep -E -v '/mocks$$')
	go tool cover -func=cover.out | tail -1
	go tool cover -html=cover.out -o=cover.html
