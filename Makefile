setup:
	go mod tidy
	go mod vendor
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	go install github.com/evilmartians/lefthook@latest
	lefthook install

lint:
	golangci-lint -c ./tools/.golangci.yml run --fix

run:
	go run ./cmd/main.go
test:
	go test ./...