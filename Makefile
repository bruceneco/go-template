setup: setup-git-hooks
	go mod tidy
	go mod vendor
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

setup-git-hooks:
	cp ./tools/pre-commit ./.git/hooks
	chmod +x .git/hooks/pre-commit

lint:
	golangci-lint -c ./tools/.golangci.yml run --fix

run:
	go run ./cmd/main.go