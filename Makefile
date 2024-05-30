.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/main cmd/ladbrokes/main.go

.PHONY: run
run: build
	docker-compose up --force-recreate --build --renew-anon-volumes

.PHONY: lint
lint:
	golangci-lint run
