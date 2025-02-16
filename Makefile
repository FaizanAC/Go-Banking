hello:
	echo "Hello"

build:
	docker build -t go-banking .

run:
	go run cmd/main.go

test:
	docker compose up --build