# Auth Service

The Auth Service is responsible for handling user authentication and authorization. It provides endpoints for user
login, and token validation (simple JWT validation).

## Prerequisites

Ensure you have the following installed:

- **GO >= 1.23.3**: [Install GO](https://go.dev/doc/install)
- **Docker**: [Install Docker](https://docs.docker.com/get-started/get-docker/)
- **Docker Compose**: [Install Docker Compose](https://docs.docker.com/compose/install/)

## Project Structure

The project structure for your system will look like this:

```text
/auth-service
├── /internal # Internal source code for the service
│   ├── app # Bootstrap code for the service 
│   ├── config # Configuration related code
│   ├── handler # HTTP handlers related code (controllers)
│   ├── repositories # Data layer code for the service 
│   ├── services # Core business logic code 
├── .env.example # Example environment variables
├── Dockerfile # Dockerfile for building the system 
├── go.mod # Go module file
├── go.sum # Go module file
├── main.go # Main entry point for the system 
├── Makefile # Makefile for building and running the system
├── README.md # This setup guide
```

## Running the System

If you have `Make` installed, you can use the `Makefile` to run the system.

```shell
  make run 
```

If you don't have Make installed, you can use `Go Command` directly:

```shell
  go run main.go
```

## API Endpoints

The Auth Service provides the following API endpoints:

- `POST /api/v1/login`: Login with username and password to get a JWT token.
- `GET /api/v1/me`: Validate a JWT token and retrieve user information.

## Environment Variables

You can find the environment variables in the `.env.example` file. You can copy this file to `.env` and update the
values.

## Accessing the Service

You can access the service at `http://0.0.0.0`.

## Create a New User

If you have `Make` installed, you can use the `Makefile` to create a new user:

```shell
  make create_user
```
or you can use `Go Command` directly:

```shell
cd cmd && go run mongo_create_user_cmd.go
```

If you are inside a Docker container, you can use the following command:

```shell
docker exec -it <container_id> ./mongo_create_user_cmd # Using Docker
docker compose exec <container_id> ./mongo_create_user_cmd # Using Docker Compose
```

## Testing

To run the tests, you can use the `Makefile`:

```shell
  make test
```

Or use `Go Command` directly:

```shell
  go test -v -cover -race ./...
```