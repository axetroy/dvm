test:
	go test --cover -covermode=count -coverprofile=coverage.out ./...

build:
	bash build.sh
