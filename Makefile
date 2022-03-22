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

update-go-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy -go=1.18
ifneq (,$(wildcard vendor))
	go mod vendor
endif