install:
	go mod tidy

build:
	go build -o bin/app 

run: build
	./bin/app

unit:
	go test -v ./...