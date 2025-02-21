# Go Banking
Go Banking is a REST Banking Service written in GoLang

## Development

This repository uses `GoDotEnv` for managing environment variables.

Create a `.env` file in the root directory with the following:

```
PORT:
HOST:
JWT_KEY:
POSTGRES_DB:
POSTGRES_USER:
POSTGRES_PASSWORD:
```

Authentication is done via JWT -  https://github.com/golang-jwt/jwt

Routing is done via Gin - https://github.com/gin-gonic/gin


## Build Docker Images

```
make build
```

## Run App Locally

```
make run
```

## Running Tests

```
make test
```
