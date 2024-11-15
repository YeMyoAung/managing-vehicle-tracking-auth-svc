run:
	@go mod tidy && go run main.go

build:
	@go mod tidy && go build -o bin/auth-svc

test:
	@go test -v -cover -race ./...

create_user:
	@go run cmd/mongo_create_user_cmd.go

build_create_user:
	@go build -o bin/mongo_create_user_cmd bin/mongo_create_user_cmd.go

.PHONY: run build test