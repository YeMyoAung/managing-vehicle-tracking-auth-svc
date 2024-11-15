# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY main.go .
COPY go.mod .
COPY go.sum .
COPY internal/ internal/
COPY cmd/mongo_create_user_cmd.go cmd/mongo_create_user_cmd.go

RUN go mod tidy
RUN go build -o bin/auth-svc
RUN go build -o bin/mongo_create_user_cmd cmd/mongo_create_user_cmd.go

# Stage 2: Set up the final image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bin/auth-svc .
COPY --from=builder /app/bin/mongo_create_user_cmd .
COPY .env.uat .env

EXPOSE 10000

CMD ["./auth-svc"]
