BINARY_NAME=analyzer
DOCKER_IMAGE=image-analyzer
PLATFORM=$(shell uname -m) // default arm64 for M1 mac

.PHONY: build
build:
	go build -o $(BINARY_NAME) -v ./cmd/zombie/main.go

.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*

.PHONY: test
test:
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: deps
deps:
	go get -u ./...
	go mod tidy

.PHONY: docker-build
docker-build:
	docker build --platform linux/$(PLATFORM) -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run:
	docker run --platform linux/$(PLATFORM) -p 8080:8080 $(DOCKER_IMAGE)

.PHONY: run
run:
	go run ./cmd/zombie/main.go

.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 -v ./cmd/main.go
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 -v ./cmd/main.go

.PHONY: lint
lint:
	golangci-lint run
