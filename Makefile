.PHONY: deps run build build-release clean test lint

deps:
	go mod download

run:
	cd deploy/local && docker-compose up --build --remove-orphans

build: clean
	go build -o ./bin/server ./cmd/service/main.go

build-release: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/server ./cmd/service/main.go

clean:
	rm -rf ./bin

test:
	go test -cover -v ./...

lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.40.0
	./bin/golangci-lint run ./...
	rm -rf ./bin