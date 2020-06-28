test:
	go test --cover -covermode=count -coverprofile=coverage.out ./...

build:
	bash build.sh

lint:
	golangci-lint run ./... -v

format:
	go fmt ./...