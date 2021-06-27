deps:
	go mod download

run:
	cd deploy/local && docker-compose up --build --remove-orphans

build:
	go build -o ./bin/server ./cmd/service/main.go

build-release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/server ./cmd/service/main.go

clean:
	rm -rf ./bin

test: # runs tests
	go test -cover -v ./...

lint: # runs linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.40.0
	./bin/golangci-lint run ./...
	rm -rf ./bin
