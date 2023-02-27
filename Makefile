install:
	go mod tidy

build:
	go build -o bin/app 

run: build
	./bin/app

unit:
	go clean -testcache && go test -v ./...