.PHONY: coverage
coverage:
	go test -covermode=count -coverprofile=cover.out ./border ./common ./converters ./image
	go tool cover -html cover.out -o cover.html

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run
