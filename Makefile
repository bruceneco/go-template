setup:
	go mod tidy
	go mod vendor
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	go install github.com/evilmartians/lefthook@latest
	lefthook install
	git config --local core.hooksPath .git/hooks
	go install github.com/rubenv/sql-migrate/...@latest

lint:
	golangci-lint -c ./tools/.golangci.yml run --fix

GO_ENV ?= development
run:
	sql-migrate up -config=./tools/db/dbconfig.yml -env=$(GO_ENV)
	GO_ENV=$(GO_ENV) go run ./cmd/main.go
test:
	go test ./...