setup:
	go mod tidy &\
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5 &\
	go install github.com/evilmartians/lefthook@latest &\
	go install github.com/rubenv/sql-migrate/...@latest &\
	wait
	lefthook install
	git config --local core.hooksPath .git/hooks

lint:
	golangci-lint -c ./tools/.golangci.yml run --fix

GO_ENV ?= development
run: migrate-up
	GO_ENV=$(GO_ENV) go run ./cmd/main.go

migrate-up:
	sql-migrate up -config=./tools/db/dbconfig.yml -env=$(GO_ENV)

migrate-down:
	sql-migrate down -config=./tools/db/dbconfig.yml -env=$(GO_ENV)

migrate-status:
	sql-migrate status -config=./tools/db/dbconfig.yml -env=$(GO_ENV)

new-migration:
	@read -p "Migration name: " NAME; \
    	if [ -z "$$NAME" ]; then \
    		echo "Invalid migration name"; \
    		exit 1; \
    	fi; \
    	sql-migrate new -config=./tools/db/dbconfig.yml $$NAME
test:
	go test ./...

protoc-generate: install-protoc-deps
	mkdir -p ./internal/adapters/grpc/proto/gen
	protoc  --proto_path=./internal/adapters/grpc/proto \
			--go_out=paths=source_relative:./internal/adapters/grpc/proto/gen \
			--go-grpc_out=paths=source_relative:./internal/adapters/grpc/proto/gen \
			./internal/adapters/grpc/proto/*.proto \
			--experimental_allow_proto3_optional
	protoc-go-inject-tag -input=./internal/adapters/grpc/proto/gen/*.pb.go
	go mod tidy

install-protoc-deps:
	go install github.com/favadi/protoc-go-inject-tag@latest &\
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest &\
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest &\
	wait
