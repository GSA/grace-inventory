default: integration_test

integration_test: lint
	go test -v tests_test.go

lint: dependencies
	golangci-lint run ./...

dependencies: golangci go.sum init

golangci:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

go.sum: go.mod
	go mod tidy

go.mod:
	go mod init

init:
	mkdir -p ../release
	touch ../release/grace-inventory-lambda.zip
