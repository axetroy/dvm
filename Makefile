test:
	go test --cover -covermode=count -coverprofile=coverage.out ./...

build:
	goreleaser release --snapshot --rm-dist --skip-publish

lint:
	golangci-lint run ./... -v

format:
	go fmt ./...

format-check:
	gofmt -l ./internal main.go